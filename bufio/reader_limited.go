package bufio

import (
	"bytes"
	"io"
)

func (b *Reader) LimitReader(n int) *LimitedReader {
	b.limited.n = n
	return b.limited
}

var _ BufioReader = (*LimitedReader)(nil)

type LimitedReader struct {
	r *Reader
	n int
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if l.n <= 0 {
		return 0, io.EOF
	}

	if len(p) > l.n {
		p = p[:l.n]
	}

	n, err = l.r.Read(p)
	l.n -= n
	return
}

func (l *LimitedReader) ReadByte() (c byte, err error) {
	if l.n <= 0 {
		return 0, io.EOF
	}

	c, err = l.r.ReadByte()

	if err == nil {
		l.n--
	}

	return
}

func (b *LimitedReader) Discard(n int) (discarded int, err error) {
	discarded, err = b.r.Discard(min(n, b.n))
	b.n -= discarded

	if err == nil && discarded < n {
		err = io.EOF
	}

	return
}

func (b *LimitedReader) DiscardUntil(c byte) (discarded int, err error) {
	for {
		if b.r.r == b.r.w {
			err = b.r.fill()
		}

		end := min(b.r.r+b.n, b.r.w)
		idx := bytes.IndexByte(b.r.buf[b.r.r:end], c)

		if idx >= 0 {
			discarded += idx
			b.r.r += idx
			b.n -= idx
			return
		}

		discarded += (end - b.r.r)
		b.r.r += (end - b.r.r)
		b.n -= (end - b.r.r)

		if err != nil {
			break
		}
	}

	return
}

func (b *LimitedReader) Peek(n int) (buf []byte, err error) {
	if n > b.n {
		buf, err = b.r.Peek(b.n)

		if err == nil {
			err = io.EOF
		}

		return
	}

	return b.r.Peek(n)
}
func (b *LimitedReader) ReadBytes(n int) (r []byte, err error) {
	r, err = b.Peek(n)
	b.r.r += len(r)
	b.n -= len(r)
	return
}

func (b *LimitedReader) ReadSlice(delimiter byte) (slice []byte, err error) {
	locked := b.Lock()
	start := b.r.r

	for {
		if b.r.r == b.r.w {
			if err = b.r.fill(); err != nil {
				break
			}
		}

		end := min(b.r.r+b.n, b.r.w)
		idx := bytes.IndexByte(b.r.buf[b.r.r:end], delimiter)

		if idx >= 0 {
			b.r.r += idx
			b.n -= idx
			break
		}

		b.r.r = end
	}

	if !locked {
		b.Unlock()
	}

	slice = b.r.buf[start:b.r.r]
	return
}

func (b *LimitedReader) Buffered() int {
	return min(b.n, b.r.Buffered())
}

func (b *LimitedReader) Lock() bool   { return b.r.Lock() }
func (b *LimitedReader) Unlock() bool { return b.r.Unlock() }
func (b *LimitedReader) Size() int    { return b.r.Size() }
