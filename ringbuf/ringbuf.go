package ringbuf

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

const (
	// BufferSize is the logical size of the ring buffer.
	BufferSize = 4096
	// SlackSize is reserved for potential slack operations.
	SlackSize = 4096
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
// It maintains two cursors:
//   - read: marks the beginning of the unread region
//   - write: marks the end of the unread region (and beginning of free space)
//
// The free region is defined as the region from write to read (cyclically).
type RingBuf struct {
	// buf is the underlying fixed storage of size TotalSize.
	buf [TotalSize]byte
	// read marks the beginning of the unread region.
	read uint64 // unread region: [read, write)
	// write marks the end of the unread region.
	write uint64 // free region: [write, read) cyclically.
}

func (rb *RingBuf) Reset() {
	rb.read, rb.write = 0, 0
}

// unread returns the number of unread bytes.
func (rb *RingBuf) unread() uint64 {
	return rb.write - rb.read
}

// free returns the number of free bytes available.
func (rb *RingBuf) free() uint64 {
	return BufferSize - (rb.write - rb.read)
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

		if err != nil {
			if err == io.ErrShortBuffer || err == io.EOF {
				err = nil
			}
			return n, err
		}

		if m == 0 {
			break
		}
	}
	return n, nil
}

// FillFrom does one (1) read from an io.Reader into the free buffer.
func (rb *RingBuf) FillFrom(r io.Reader) (n int64, err error) {
	free := BufferSize - (rb.write - rb.read)
	if free == 0 {
		return 0, io.ErrShortBuffer
	}

	// Compute the contiguous free region starting from the current write pointer.
	index := rb.write & ringMask
	contig := BufferSize - index
	if free < contig {
		contig = free
	}

	m, err := r.Read(rb.buf[index : index+contig])
	if m > 0 {
		rb.write += uint64(m)
		n = int64(m)
	}
	return n, err
}

// ReadByte implements io.ByteReader.
func (rb *RingBuf) ReadByte() (byte, error) {
	if rb.unread() == 0 {
		return 0, io.EOF
	}

	index := rb.read & ringMask
	b := rb.buf[index]
	rb.read++
	return b, nil
}

// Peek returns the next n unread bytes from the ring buffer without advancing the read pointer.
// If the requested data is contiguous, it returns a direct slice into rb.buf.
// If the data wraps around, it copies the wrapped portion into the slack space and returns a contiguous slice.
func (rb *RingBuf) Peek(n int) (buf []byte, err error) {
	nn := uint64(n)

	// Check that we aren't requesting more bytes than available.
	if nn > rb.unread() {
		return nil, io.EOF
	}

	start := rb.read & ringMask
	end := (rb.read + nn) & ringMask

	if start > end {
		copy(rb.buf[BufferSize:], rb.buf[:end])
	}

	return rb.buf[start : start+nn], nil
}

func (b *RingBuf) ReadBytes(n int) (r []byte, err error) {
	r, err = b.Peek(n)
	b.read += uint64(len(r))
	return
}

func (rb *RingBuf) DebugDump(wr io.Writer) {
	w := bufio.NewWriter(wr)

	// Compute the current physical positions for read and write within the main buffer.
	readPos := rb.read & ringMask
	writePos := rb.write & ringMask

	// Iterate over the full underlying array (main buffer + slack area).
	for i := uint64(0); i < TotalSize; i++ {
		b := rb.buf[i]
		// Determine ASCII representation: printable characters are shown, others as '.'
		var ch rune
		if unicode.IsPrint(rune(b)) {
			ch = rune(b)
		} else {
			ch = '.'
		}

		// Build the line with hex and ASCII.
		line := fmt.Sprintf("%02X %c", b, ch)

		// If we're in the main buffer area, add markers for read and write positions.
		if i < BufferSize {
			isRead := (i == readPos)
			isWrite := (i == writePos)
			if isRead && isWrite {
				line += " <- R + W"
			} else if isRead {
				line += " <- R"
			} else if isWrite {
				line += " <- W"
			}
		} else if i == BufferSize {
			// Mark the beginning of the slack area.
			line += " <- SLACK BEGIN"
		}

		fmt.Fprintln(w, line)
	}
	w.Flush()
}
