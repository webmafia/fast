package binary

import (
	"bufio"
	"io"

	"github.com/webmafia/fast"
)

var _ Reader = (*StreamReader)(nil)

// A StreamReader is used to efficiently read binary data from
// an io.Reader.
type StreamReader struct {
	buf *bufio.Reader
	err error
}

func NewStreamReader(r io.Reader) *StreamReader {
	return &StreamReader{
		buf: bufio.NewReader(r),
	}
}

func (b *StreamReader) Error() error {
	return b.err
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *StreamReader) Len() int {
	return b.buf.Buffered()
}

// Cap returns the capacity of the StreamReader's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *StreamReader) Cap() int {
	return b.buf.Size()
}

// Reset resets the StreamReader to be empty.
func (b *StreamReader) Reset(r io.Reader) {
	b.buf.Reset(r)
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *StreamReader) Read(dst []byte) (n int, err error) {
	return b.buf.Read(dst)
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *StreamReader) ReadByte() (byte, error) {
	return b.buf.ReadByte()
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *StreamReader) ReadBytes(n int) []byte {
	buf := make([]byte, n)
	n, _ = b.buf.Read(buf)
	return buf[:n]
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *StreamReader) ReadString(n int) string {
	return fast.BytesToString(b.ReadBytes(n))
}

func (b *StreamReader) WriteTo(w io.Writer) (n int64, err error) {
	return b.buf.WriteTo(w)
}
