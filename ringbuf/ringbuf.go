package ringbuf

import (
	"errors"
	"io"
)

const (
	// BufferSize is the logical size of the ring buffer.
	BufferSize = 4096
	// SlackSize is reserved for potential slack operations.
	SlackSize = 2048
	// TotalSize is the total size of the underlying array.
	TotalSize = BufferSize + SlackSize

	ringMask uint64 = BufferSize - 1
)

var (
	_ io.ReadWriter = (*RingBuf)(nil)
	_ io.ByteReader = (*RingBuf)(nil)
	_ io.ReaderFrom = (*RingBuf)(nil)
)

// RingBuf is a fixed-size ring buffer that implements io.Reader and io.Writer.
// It maintains three cursors:
//   - start: marks the beginning of the locked region (data that has been read but not flushed)
//   - read: marks the beginning of the unread region
//   - write: marks the end of the unread region (and beginning of free space)
//
// The free region is defined as the region from write up to start (cyclically).
type RingBuf struct {
	// buf is the underlying fixed storage of size TotalSize.
	buf [TotalSize]byte
	// start marks the beginning of the locked region.
	start uint64 // locked region: [start, read)
	// read marks the beginning of the unread region.
	read uint64 // unread region: [read, write)
	// write marks the end of the unread region.
	write uint64 // free region: [write, start)
}

// Use plain arithmetic for the differences.
func (rb *RingBuf) locked() uint64 {
	return rb.read - rb.start
}

func (rb *RingBuf) unread() uint64 {
	return rb.write - rb.read
}

func (rb *RingBuf) free() uint64 {
	return BufferSize - (rb.read - rb.start) - (rb.write - rb.read)
}

func (rb *RingBuf) Flush() {
	rb.start = rb.read
}

// Read implements io.Reader.
func (rb *RingBuf) Read(p []byte) (n int, err error) {
	avail := rb.unread()
	if avail == 0 {
		return 0, io.EOF
	}

	toRead := uint64(len(p))
	if toRead > avail {
		toRead = avail
	}

	// First part: from rb.read (mod BufferSize) up to the end of the logical ring.
	start := rb.read & ringMask
	first := BufferSize - start
	if toRead < first {
		first = toRead
	}
	n1 := copy(p, rb.buf[start:start+first])
	n += n1

	// If there is wrap-around, copy the remainder from the beginning.
	remaining := int(toRead) - n1
	if remaining > 0 {
		n2 := copy(p[n1:], rb.buf[:remaining])
		n += n2
	}

	// Advance the read pointer.
	rb.read += uint64(n)

	// Flush the read part directly
	rb.Flush()

	return n, nil
}

// Write implements io.Writer.
func (rb *RingBuf) Write(p []byte) (n int, err error) {
	free := rb.free()
	if free == 0 {
		return 0, io.ErrShortBuffer
	}

	toWrite := uint64(len(p))
	if toWrite > free {
		toWrite = free
	}

	// First part: from rb.write (mod BufferSize) to the end of the logical ring.
	start := rb.write & ringMask
	first := BufferSize - start
	if toWrite < first {
		first = toWrite
	}
	n1 := copy(rb.buf[start:start+first], p)
	n += n1

	// If wrap-around is needed, copy the remainder to the beginning.
	remaining := int(toWrite) - n1
	if remaining > 0 {
		n2 := copy(rb.buf[:remaining], p[n1:])
		n += n2
	}

	// Advance the write pointer.
	rb.write += uint64(n)

	return n, nil
}

// ReadFrom implements io.ReaderFrom.
// It repeatedly calls FillFrom until an error occurs or the buffer is full.
func (rb *RingBuf) ReadFrom(r io.Reader) (n int64, err error) {
	for {
		m, err := rb.FillFrom(r)
		n += m

		// If we got a non-nil error, then:
		// - If it's io.ErrShortBuffer, the ring buffer is full; we return what we've read so far.
		// - If it's io.EOF, we return the data and a nil error.
		// - Otherwise, we return the error.
		if err != nil {
			if err == io.ErrShortBuffer || err == io.EOF {
				err = nil
			}

			return n, err
		}

		// If no bytes were read, break out to avoid an infinite loop.
		if m == 0 {
			break
		}
	}
	return n, nil
}

// FillFrom does one (1) read from an io.Reader into the free buffer.
func (rb *RingBuf) FillFrom(r io.Reader) (n int64, err error) {
	// Calculate free space available.
	free := BufferSize - (rb.write - rb.start)
	if free == 0 {
		return 0, io.ErrShortBuffer
	}

	// Compute the contiguous free region starting from the current write pointer.
	index := rb.write & ringMask
	contig := BufferSize - index
	if free < contig {
		contig = free
	}

	// Read directly into the contiguous free region.
	m, err := r.Read(rb.buf[index : index+contig])
	if m > 0 {
		rb.write += uint64(m)
		n = int64(m)
	}
	return n, err
}

// ReadByte implements io.ByteReader.
func (rb *RingBuf) ReadByte() (byte, error) {
	// If there's no unread data, return EOF.
	if rb.unread() == 0 {
		return 0, io.EOF
	}

	// Determine the index in the underlying buffer.
	index := rb.read & ringMask
	b := rb.buf[index]

	// Advance the read pointer.
	rb.read++

	// Flush the read portion, releasing the locked data.
	rb.Flush()

	return b, nil
}

// Peek returns the next n unread bytes from the ring buffer without advancing the read pointer.
// If the requested data is contiguous, it returns a direct slice into rb.buf.
// If the data wraps around, it copies the wrapped portion into the slack space and returns a contiguous slice.
func (rb *RingBuf) Peek(n int) (buf []byte, err error) {
	// Ensure there is enough unread data.
	if uint64(n) > rb.unread() {
		return nil, io.EOF
	}

	// Compute the starting index in the underlying buffer.
	start := int(rb.read & ringMask)
	contig := int(BufferSize) - start

	if n <= contig {
		// Data is contiguous.
		return rb.buf[start : start+n], nil
	}

	// Data wraps around: calculate how many bytes wrap.
	remainder := n - contig
	// Ensure the slack space is sufficient.
	if remainder > int(SlackSize) {
		return nil, errors.New("ringbuf: insufficient slack space for Peek")
	}

	// Copy the wrapped portion (from the beginning of rb.buf) into the slack area.
	copy(rb.buf[BufferSize:BufferSize+remainder], rb.buf[:remainder])

	// Return a contiguous slice from the read index through the slack area.
	return rb.buf[start : BufferSize+remainder], nil
}
