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
	_ io.ByteReader   = (*Buffer)(nil)
)

// A byte buffer, highly optimized to minimize allocations and GC pressure.
type Buffer struct {
	B   []byte
	pos int
	r   io.Reader
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
	b.ResetReader(nil)
}

func (b *Buffer) ResetReader(r io.Reader) {
	b.pos = 0
	b.r = r
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

// Release resets the reader state after processing a message.
// It discards consumed data and prepares for the next message.
func (b *Buffer) Release() {
	b.release(0)
}

func (b *Buffer) ReleaseBefore(n int) {
	b.release(b.croppedPos(n))
}

func (b *Buffer) ReleaseAfter(n int) {
	b.pos = b.croppedPos(n)
	b.B = b.B[:b.pos]
}

func (b *Buffer) croppedPos(n int) int {
	return max(0, min(n, len(b.B)))
}

//go:inline
func (b *Buffer) release(n int) {
	b.B = append(b.B[n:n], b.B[b.pos:]...)
	b.pos = n
}

func (b *Buffer) Pos() int {
	return b.pos
}

func (b *Buffer) SetPos(n int) {
	b.pos = b.croppedPos(n)
}

func (b *Buffer) AdjustPos(delta int) {
	b.pos = b.croppedPos(b.pos + delta)
}

func (b *Buffer) Fill(n int) error {
	required := n - (len(b.B) - b.pos)

	if required <= 0 {
		return nil
	}

	return b.fillFromReader(required)
}

// fill ensures that the buffer contains at least `n` bytes from the current position.
func (b *Buffer) fillFromReader(required int) (err error) {
	if b.r == nil {
		return io.EOF
	}

	if err = b.Grow(required); err != nil {
		return
	}

	size := len(b.B)
	b.B = b.B[:cap(b.B)]

	// Read data into the buffer.
	nRead, err := io.ReadAtLeast(b.r, b.B[size:], required)
	b.B = b.B[:size+nRead]

	return
}

// consume advances the read position by `n` bytes.
func (b *Buffer) consume(n int) {
	b.pos += n
}

// ReadByte implements io.ByteReader.
func (b *Buffer) ReadByte() (c byte, err error) {
	c, err = b.PeekByte()

	if err == nil {
		b.consume(1)
	}

	return
}

// ReadByte implements io.ByteReader.
func (b *Buffer) PeekByte() (c byte, err error) {
	if err = b.Fill(1); err != nil {
		return
	}

	return b.B[b.pos], nil
}

func (b *Buffer) ReadBytes(n int) (bytes []byte, err error) {
	bytes, err = b.PeekBytes(n)

	if err == nil {
		b.consume(n)
	}

	return
}

func (b *Buffer) PeekBytes(n int) (bytes []byte, err error) {
	if err = b.Fill(n); err != nil {
		return
	}

	return b.B[b.pos : b.pos+n], nil
}

func (b *Buffer) ReadString(n int) (str string, err error) {
	bytes, err := b.PeekBytes(n)

	if err == nil {
		str = fast.BytesToString(bytes)
		b.consume(n)
	}

	return
}
