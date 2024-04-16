package fast

import (
	"unicode/utf8"
)

var _ Writer = (*BinaryBuffer)(nil)

// A BinaryBuffer is used to efficiently build a string using Write methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero BinaryBuffer.
type BinaryBuffer struct {
	buf []byte
}

func NewBinaryBuffer(cap int) *BinaryBuffer {
	return &BinaryBuffer{
		buf: make([]byte, 0, cap),
	}
}

// String returns the accumulated string.
func (b BinaryBuffer) String() string {
	return BytesToString(b.buf)
}

// String returns the accumulated string as bytes.
func (b BinaryBuffer) Bytes() []byte {
	return b.buf
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *BinaryBuffer) Len() int {
	return len(b.buf)
}

// Cap returns the capacity of the BinaryBuffer's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *BinaryBuffer) Cap() int {
	return cap(b.buf)
}

// Reset resets the BinaryBuffer to be empty.
func (b *BinaryBuffer) Reset() {
	b.buf = b.buf[:0]
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *BinaryBuffer) grow(n int) {
	buf := MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
	copy(buf, b.buf)
	b.buf = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *BinaryBuffer) Grow(n int) {
	if n < 0 {
		panic("fast.BinaryBuffer.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *BinaryBuffer) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *BinaryBuffer) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *BinaryBuffer) WriteRune(r rune) (int, error) {
	n := len(b.buf)
	b.buf = utf8.AppendRune(b.buf, r)
	return len(b.buf) - n, nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *BinaryBuffer) WriteString(s string) (int, error) {
	b.buf = append(b.buf, s...)
	return len(s), nil
}
