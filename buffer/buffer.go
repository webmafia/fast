package buffer

import (
	"fmt"
	"io"

	"github.com/webmafia/fast"
)

var (
	_ io.Writer       = (*Buffer)(nil)
	_ io.ByteWriter   = (*Buffer)(nil)
	_ io.StringWriter = (*Buffer)(nil)
	_ io.ReaderFrom   = (*Buffer)(nil)
	_ io.WriterTo     = (*Buffer)(nil)
	_ fmt.Stringer    = (*Buffer)(nil)
)

// A byte buffer, highly optimized to minimize allocations and GC pressure.
type Buffer struct {
	B []byte
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		B: fast.MakeNoZeroCap(0, size),
	}
}

func (b *Buffer) Len() int {
	return len(b.B)
}

func (b *Buffer) Cap() int {
	return cap(b.B)
}

func (b *Buffer) Reset() {
	b.B = b.B[:0]
}

func (b *Buffer) Bytes() []byte {
	return b.B
}

// WriteString implements fmt.Stringer.
func (b *Buffer) String() string {
	return fast.BytesToString(b.B)
}

// Write implements io.Writer.
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

// WriteByte implements io.ByteWriter.
func (b *Buffer) WriteByte(c byte) error {
	b.B = append(b.B, c)
	return nil
}

// WriteString implements io.StringWriter.
func (b *Buffer) WriteString(s string) (n int, err error) {
	b.B = append(b.B, s...)
	return len(s), nil
}

// ReadFrom implements io.ReaderFrom.
func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.B
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart

	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}

	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)

		if err != nil {
			b.B = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)

	if err == nil {
		b.Reset()
	}

	return int64(n), err
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *Buffer) grow(n int) {
	buf := fast.MakeNoZero(2*cap(b.B) + n)[:len(b.B)]
	copy(buf, b.B)
	b.B = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *Buffer) Grow(n int) error {
	if n < 0 {
		return ErrNegativeCount
	}

	if cap(b.B)-len(b.B) < n {
		b.grow(n)
	}

	return nil
}
