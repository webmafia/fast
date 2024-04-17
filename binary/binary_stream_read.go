package fast

import (
	"bufio"
	"io"

	"github.com/webmafia/fast"
)

var _ Reader = (*BinaryStreamReader)(nil)

// A BinaryStreamReader is used to efficiently read binary data from
// an io.Reader.
type BinaryStreamReader struct {
	buf *bufio.Reader
	err error
}

func NewBinaryStreamReader(r io.Reader) *BinaryStreamReader {
	return &BinaryStreamReader{
		buf: bufio.NewReader(r),
	}
}

func (b *BinaryStreamReader) Error() error {
	return b.err
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *BinaryStreamReader) Len() int {
	return b.buf.Buffered()
}

// Cap returns the capacity of the BinaryStreamReader's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *BinaryStreamReader) Cap() int {
	return b.buf.Size()
}

// Reset resets the BinaryStreamReader to be empty.
func (b *BinaryStreamReader) Reset(r io.Reader) {
	b.buf.Reset(r)
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *BinaryStreamReader) Read(dst []byte) (n int, err error) {
	return b.buf.Read(dst)
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *BinaryStreamReader) ReadByte() (byte, error) {
	return b.buf.ReadByte()
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b *BinaryStreamReader) ReadRune() (r rune, size int, err error) {
	return b.buf.ReadRune()
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *BinaryStreamReader) ReadBytes(n int) []byte {
	buf := make([]byte, n)
	n, _ = b.buf.Read(buf)
	return buf[:n]
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *BinaryStreamReader) ReadString(n int) string {
	return fast.BytesToString(b.ReadBytes(n))
}

func (b *BinaryStreamReader) WriteTo(w io.Writer) (n int64, err error) {
	return b.buf.WriteTo(w)
}
