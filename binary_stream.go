package fast

import (
	"bufio"
	"io"
)

var _ Writer = BinaryStream{}

// A BinaryStream is used to efficiently write binary data to
// an io.Writer.
type BinaryStream struct {
	buf *bufio.Writer
}

func NewBinaryStream(w io.Writer) BinaryStream {
	return BinaryStream{
		buf: bufio.NewWriter(w),
	}
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b BinaryStream) Len() int {
	return b.buf.Buffered()
}

// Cap returns the capacity of the BinaryStream's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b BinaryStream) Cap() int {
	return b.buf.Size()
}

// Reset resets the BinaryStream to be empty.
func (b BinaryStream) Reset(w io.Writer) {
	b.buf.Reset(w)
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b BinaryStream) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b BinaryStream) WriteByte(c byte) error {
	return b.buf.WriteByte(c)
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to b's buffer.
// It returns the length of r and a nil error.
func (b BinaryStream) WriteRune(r rune) (int, error) {
	return b.buf.WriteRune(r)
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b BinaryStream) WriteString(s string) (int, error) {
	return b.buf.WriteString(s)
}

func (b BinaryStream) Flush() error {
	return b.buf.Flush()
}

func (b BinaryStream) ReadFrom(r io.Reader) (n int64, err error) {
	return b.buf.ReadFrom(r)
}
