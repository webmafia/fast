package ringbuf

import (
	"bytes"
	"io"
)

type RingBufferReader interface {
	Buffered() int
	Reset(rd io.Reader)
	ResetBytes(b []byte)
	Read(p []byte) (n int, err error)
	ReadByte() (byte, error)
	Peek(n int) ([]byte, error)
	ReadBytes(n int) (b []byte, err error)
	Discard(n int) (discarded int, err error)
	DiscardUntil(c byte) (discarded int, err error)
}

var _ RingBufferReader = (*Reader)(nil)

// Reader wraps an io.Reader and uses a RingBuf for buffering.
type Reader struct {
	r       io.Reader
	ring    RingBuf
	limited *LimitedReader
	err     error // sticky error from the underlying reader, if any
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

func (r *Reader) Reset(rd io.Reader) {
	if rd == r {
		return
	}

	r.r = rd
	r.ring.Reset()
}

func (r *Reader) ResetBytes(b []byte) {
	if br, ok := r.r.(*bytes.Reader); ok {
		br.Reset(b)
	} else {
		r.r = bytes.NewReader(b)
	}

	r.ring.Reset()
}

func (r *Reader) SetManualFlush(v bool) {
	r.ring.manualFlush = v
}

func (r *Reader) Flush() {
	r.ring.Flush()
}

// fill attempts to read from the underlying reader into the ring buffer.
// It uses FillFrom and stores any non-EOF error into r.err.
func (r *Reader) fill() {
	// Only fill if there's free space.
	if r.ring.free() == 0 || r.err != nil {
		return
	}
	// FillFrom will fill as much as possible.
	_, err := r.ring.FillFrom(r.r)

	// If err == io.EOF, we do not mark it as sticky,
	// so that the buffered data can still be read.
	if err != nil && err != io.EOF {
		r.err = err
	}
}

// Returns the total number of read bytes from underlying io.Reader since last reset.
func (r *Reader) TotalRead() int {
	return int(r.ring.write)
}

// Buffered returns the number of unread bytes currently buffered.
func (r *Reader) Buffered() int {
	return int(r.ring.unread())
}

// Read implements io.Reader.
// It first ensures that there is data in the ring buffer by calling fill() as needed.
func (r *Reader) Read(p []byte) (n int, err error) {
	total := 0
	for len(p) > 0 {
		// If no buffered data, attempt to fill.
		if r.ring.unread() == 0 {
			r.fill()
			// If still empty, return error if no data read.
			if r.ring.unread() == 0 {
				if total > 0 {
					return total, nil
				}
				if r.err != nil {
					return 0, r.err
				}
				return 0, io.EOF
			}
		}
		nn, _ := r.ring.Read(p)
		total += nn
		p = p[nn:]
		// If we read less than requested and there is a sticky error, break.
		if len(p) > 0 && r.err != nil {
			break
		}
	}
	return total, nil
}

// ReadByte implements io.ByteReader.
// It refills the ring buffer if necessary.
func (r *Reader) ReadByte() (byte, error) {
	if r.ring.unread() == 0 {
		r.fill()
		if r.ring.unread() == 0 {
			if r.err != nil {
				return 0, r.err
			}
			return 0, io.EOF
		}
	}
	return r.ring.ReadByte()
}

// Peek returns the next n unread bytes without advancing the read pointer.
// It refills the ring buffer until at least n bytes are available or an error occurs.
func (r *Reader) Peek(n int) ([]byte, error) {
	for uint64(n) > r.ring.unread() && r.err == nil {
		before := r.ring.unread()
		r.fill()
		// Only break if fill() didn't add any new bytes.
		if r.ring.unread() == before {
			break
		}
	}
	if uint64(n) > r.ring.unread() {
		return nil, io.EOF
	}
	return r.ring.Peek(n)
}

// ReadBytes reads exactly n bytes from the buffered Reader.
// It fills the buffer as needed; if fewer than n bytes are available, it returns io.EOF.
func (r *Reader) ReadBytes(n int) (b []byte, err error) {
	b, err = r.Peek(n)
	r.ring.read += uint64(len(b))
	return
}

// Discard skips the next n bytes, returning the number of bytes discarded.
// It fills the buffer as needed.
func (r *Reader) Discard(n int) (discarded int, err error) {
	total := 0
	for n > 0 {
		// Fill the buffer if necessary.
		if r.ring.unread() == 0 {
			r.fill()
			if r.ring.unread() == 0 {
				// No more data available.
				if total > 0 {
					return total, nil
				}
				if r.err != nil {
					return total, r.err
				}
				return total, io.EOF
			}
		}
		avail := int(r.ring.unread())
		toDiscard := n
		if toDiscard > avail {
			toDiscard = avail
		}
		r.ring.read += uint64(toDiscard)
		total += toDiscard
		n -= toDiscard
	}
	return total, nil
}

// DiscardUntil discards all bytes until (but not including) the first occurrence of c.
// It returns the number of bytes discarded. If c is not found and no more data is available,
// it returns io.EOF.
func (r *Reader) DiscardUntil(c byte) (discarded int, err error) {
	total := 0
	for {
		// Fill the buffer if needed.
		if r.ring.unread() == 0 {
			r.fill()
			if r.ring.unread() == 0 {
				if total > 0 {
					return total, io.EOF
				}
				if r.err != nil {
					return total, r.err
				}
				return total, io.EOF
			}
		}
		avail := int(r.ring.unread())
		buf, err := r.ring.Peek(avail)
		if err != nil {
			return total, err
		}
		// Use bytes.IndexByte to search for c.
		index := bytes.IndexByte(buf, c)
		if index >= 0 {
			// Found c: discard up to that position.
			r.ring.read += uint64(index)
			total += index
			return total, nil
		}
		// c not found in current buffer: discard all and continue.
		r.ring.read += uint64(avail)
		total += avail
	}
}

// LimitReader returns a pointer to a LimitedReader wrapping the current Reader.
// It updates the limit if a LimitedReader already exists.
func (r *Reader) LimitReader(n int) *LimitedReader {
	if r.limited == nil {
		r.limited = &LimitedReader{
			r: r,
			n: n,
		}
	} else {
		r.limited.n = n
	}
	return r.limited
}

func (r *Reader) DebugDump(w io.Writer) {
	r.ring.DebugDump(w)
}
