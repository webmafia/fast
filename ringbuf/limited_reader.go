package ringbuf

import (
	"io"
)

// LimitedReader wraps a pointer to a Reader and limits the total number of bytes
// that can be read. The underlying Reader and the remaining limit are stored in
// unexported fields.
type LimitedReader struct {
	r *Reader // underlying Reader (unexported)
	n int     // maximum remaining bytes allowed to be read (unexported)
}

// Buffered returns the number of unread bytes currently buffered, capped to the remaining limit.
func (lr *LimitedReader) Buffered() int {
	buf := lr.r.Buffered()
	if buf > lr.n {
		return lr.n
	}
	return buf
}

// Read implements io.Reader, reading at most lr.n bytes.
func (lr *LimitedReader) Read(p []byte) (n int, err error) {
	if lr.n <= 0 {
		return 0, io.EOF
	}
	// Ensure we don't read more than the remaining limit.
	if len(p) > lr.n {
		p = p[:lr.n]
	}
	n, err = lr.r.Read(p)
	lr.n -= n
	return n, err
}

// ReadByte implements io.ByteReader.
func (lr *LimitedReader) ReadByte() (byte, error) {
	if lr.n <= 0 {
		return 0, io.EOF
	}
	b, err := lr.r.ReadByte()
	if err == nil {
		lr.n--
	}
	return b, err
}

// Peek returns the next n unread bytes without advancing the read pointer,
// but no more than the remaining limit.
func (lr *LimitedReader) Peek(n int) ([]byte, error) {
	if lr.n <= 0 {
		return nil, io.EOF
	}
	if n > lr.n {
		n = lr.n
	}
	return lr.r.Peek(n)
}

// ReadBytes reads exactly n bytes from the buffered LimitedReader.
// If fewer than n bytes are available before reaching the limit, it returns io.EOF.
func (lr *LimitedReader) ReadBytes(n int) ([]byte, error) {
	if lr.n <= 0 {
		return nil, io.EOF
	}
	if n > lr.n {
		n = lr.n
	}
	b, err := lr.r.ReadBytes(n)
	lr.n -= len(b)
	return b, err
}

// Discard skips the next n bytes, returning the number of bytes discarded.
func (lr *LimitedReader) Discard(n int) (discarded int, err error) {
	if lr.n <= 0 {
		return 0, io.EOF
	}
	if n > lr.n {
		n = lr.n
	}
	discarded, err = lr.r.Discard(n)
	lr.n -= discarded
	return discarded, err
}

// DiscardUntil discards bytes until (but not including) the first occurrence of c.
// It returns the number of bytes discarded. If c is not found and no more data is available,
// it returns io.EOF.
func (lr *LimitedReader) DiscardUntil(c byte) (discarded int, err error) {
	if lr.n <= 0 {
		return 0, io.EOF
	}
	discarded, err = lr.r.DiscardUntil(c)
	if discarded > lr.n {
		discarded = lr.n
	}
	lr.n -= discarded
	return discarded, err
}
