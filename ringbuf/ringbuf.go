package ringbuf

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
	buf   [TotalSize]byte
	start uint64 // locked region: [start, read)
	// read marks the beginning of the unread region.
	read uint64 // unread region: [read, write)
	// write marks the end of the unread region.
	write       uint64 // free region: [write, read) cyclically.
	manualFlush bool
}

func (rb *RingBuf) SetManualFlush(v bool) {
	rb.manualFlush = v

	if !v {
		rb.start = rb.read
	}
}

func (rb *RingBuf) Reset() {
	rb.start, rb.read, rb.write = 0, 0, 0
}

// unread returns the number of unread bytes.
func (rb *RingBuf) unread() uint64 {
	return rb.write - rb.read
}

// free returns the number of free bytes available.
func (rb *RingBuf) free() uint64 {
	return BufferSize - (rb.write - rb.start)
}

func (rb *RingBuf) advance(n uint64) {
	rb.read += n

	if !rb.manualFlush {
		rb.start = rb.read
	}
}

// Flushes the read region, so that it can be overwritten.
func (rb *RingBuf) Flush() {
	rb.start = rb.read
}

// Rewinds the read data back to start.
func (rb *RingBuf) Rewind() {
	rb.read = rb.start
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
	rb.advance(uint64(n))
	return n, nil
}

// Write implements io.Writer.
func (rb *RingBuf) Write(p []byte) (n int, err error) {
	free := int(rb.free())

	if free == 0 {
		return 0, io.ErrShortBuffer
	}

	n = min(len(p), free)
	buf := rb.buf[:BufferSize]

	// First part: from rb.write (mod BufferSize) to the end of the logical ring.
	start := rb.write & ringMask
	end := (rb.write + uint64(n)) & ringMask

	written := copy(buf[start:], p)

	if written < n {
		copy(buf[:end], p[written:])
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
	free := rb.free()

	if free == 0 {
		return 0, io.ErrShortBuffer
	}

	buf := rb.buf[:BufferSize]
	start := rb.write & ringMask
	end := (rb.write + free) & ringMask
	buf = buf[start:]

	if start < end {
		buf = buf[:end-start]
	}

	m, err := r.Read(buf)

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
	rb.advance(1)
	return b, nil
}

// Peek returns the next n unread bytes from the ring buffer without advancing the read pointer.
// If the requested data is contiguous, it returns a direct slice into rb.buf.
// If the data wraps around, it copies the wrapped portion into the slack space and returns a contiguous slice.
func (rb *RingBuf) Peek(n int) (buf []byte, err error) {
	if n > BufferSize {
		return nil, fmt.Errorf("tried to peek %d bytes, which exceeds the limit of %d", n, BufferSize)
	}

	nn := uint64(n)

	// Check that we aren't requesting more bytes than available.
	if nn > rb.unread() {
		return nil, io.EOF
	}

	start := rb.read & ringMask
	end := (rb.read + nn) & ringMask

	// If the data wraps around the ring buffer
	if start > end {
		copy(rb.buf[BufferSize:], rb.buf[:end])
	}

	return rb.buf[start : start+nn], nil
}

func (rb *RingBuf) ReadBytes(n int) (r []byte, err error) {
	r, err = rb.Peek(n)
	rb.advance(uint64(len(r)))
	return
}

func (rb *RingBuf) DebugDump(wr io.Writer) {
	w := bufio.NewWriter(wr)

	// Compute the current physical positions for read and write within the main buffer.
	startPos := rb.start & ringMask
	readPos := rb.read & ringMask
	writePos := rb.write & ringMask
	annot := make([]string, 0, 3)

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
			annot = annot[:0]

			if i == startPos {
				annot = append(annot, "S")
			}

			if i == readPos {
				annot = append(annot, "R")
			}

			if i == writePos {
				annot = append(annot, "W")
			}

			if len(annot) > 0 {
				line += " <- " + strings.Join(annot, " + ")
			}
		} else if i == BufferSize {
			// Mark the beginning of the slack area.
			line += " <- SLACK BEGIN"
		}

		fmt.Fprintln(w, line)
	}
	w.Flush()
}
