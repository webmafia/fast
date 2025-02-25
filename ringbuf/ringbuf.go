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

var _ io.ReadWriter = (*RingBuf)(nil)

// RingBuffer is a fixed-size ring buffer that implements io.Reader and io.Writer.
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

func (rb *RingBuf) locked() uint64 {
	return (rb.read - rb.start) & ringMask
}

func (rb *RingBuf) unread() uint64 {
	return (rb.write - rb.read) & ringMask
}

func (rb *RingBuf) free() uint64 {
	return BufferSize - rb.locked() - rb.unread()
}

// Read implements io.Reader.
func (rb *RingBuf) Read(p []byte) (n int, err error) {
	avail := rb.unread()
	if avail == 0 {
		// No unread data: standard behavior is to return io.EOF.
		return 0, io.EOF
	}

	// We only copy as many bytes as are available.
	toRead := uint64(len(p))
	if toRead > avail {
		toRead = avail
	}

	// First part: from rb.read up to the end of the logical ring.
	// Logical ring size is BufferSize.
	start := rb.read & ringMask
	first := BufferSize - start
	if toRead < first {
		first = toRead
	}
	n1 := copy(p, rb.buf[start:start+first])
	n += n1

	// If there is wrap-around, copy the remainder from the beginning of rb.buf.
	remaining := int(toRead) - n1
	if remaining > 0 {
		n2 := copy(p[n1:], rb.buf[:remaining])
		n += n2
	}

	// Advance read pointer.
	rb.read += uint64(n)

	return n, nil
}

// Write implements io.Writer.
func (rb *RingBuf) Write(p []byte) (n int, err error) {
	free := rb.free()
	if free == 0 {
		// No free space: in a non-blocking ring buffer you may return an error.
		// Here we return io.ErrShortBuffer.
		return 0, io.ErrShortBuffer
	}

	toWrite := uint64(len(p))
	if toWrite > free {
		toWrite = free
	}

	// First part: from rb.write up to the end of the logical ring.
	start := rb.write & ringMask
	first := BufferSize - start
	if toWrite < first {
		first = toWrite
	}
	n1 := copy(rb.buf[start:start+first], p)
	n += n1

	// If there is wrap-around, copy the remainder to the beginning of rb.buf.
	remaining := int(toWrite) - n1
	if remaining > 0 {
		n2 := copy(rb.buf[:remaining], p[n1:])
		n += n2
	}

	// Advance write pointer.
	rb.write += uint64(n)

	return n, nil
}
