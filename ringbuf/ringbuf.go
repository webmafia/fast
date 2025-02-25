package ringbuf

import "io"

const (
	// BufferSize is the logical size of the ring buffer.
	BufferSize uint64 = 4096
	// SlackSize is reserved for potential slack operations.
	SlackSize uint64 = 2048
	// TotalSize is the total size of the underlying array.
	TotalSize uint64 = BufferSize + SlackSize

	ringMask uint64 = BufferSize - 1
)

var (
	_ io.ReadWriter = (*RingBuf)(nil)
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
func (rb *RingBuf) ReadFrom(r io.Reader) (n int64, err error) {
	panic("unimplemented")
}
