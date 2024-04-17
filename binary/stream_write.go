package binary

import (
	"bufio"
	"io"
)

var _ Writer = StreamWriter{}

// A Stream is used to efficiently write binary data to
// an io.Writer.
type StreamWriter struct {
	buf *bufio.Writer
}

func NewStreamWriter(w io.Writer) StreamWriter {
	return StreamWriter{
		buf: bufio.NewWriter(w),
	}
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b StreamWriter) Len() int {
	return b.buf.Buffered()
}

// Cap returns the capacity of the Stream's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b StreamWriter) Cap() int {
	return b.buf.Size()
}

// Reset resets the Stream to be empty.
func (b StreamWriter) Reset(w io.Writer) {
	b.buf.Reset(w)
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b StreamWriter) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b StreamWriter) WriteByte(c byte) error {
	return b.buf.WriteByte(c)
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b StreamWriter) WriteString(s string) (int, error) {
	return b.buf.WriteString(s)
}

func (b StreamWriter) Flush() error {
	return b.buf.Flush()
}

func (b StreamWriter) ReadFrom(r io.Reader) (n int64, err error) {
	return b.buf.ReadFrom(r)
}

func (b StreamWriter) WriteVal(val any) {
	switch v := val.(type) {

	case Encoder:
		b.WriteEnc(v)

	case string:
		b.WriteString(v)

	case []byte:
		b.Write(v)

	case int:
		b.WriteInt(v)

	case int8:
		b.WriteInt8(v)

	case int16:
		b.WriteInt16(v)

	case int32:
		b.WriteInt32(v)

	case int64:
		b.WriteInt64(v)

	case uint:
		b.WriteUint(v)

	case uint8:
		b.WriteUint8(v)

	case uint16:
		b.WriteUint16(v)

	case uint32:
		b.WriteUint32(v)

	case uint64:
		b.WriteUint64(v)

	case float32:
		b.WriteFloat32(v)

	case float64:
		b.WriteFloat64(v)

	case bool:
		b.WriteBool(v)

	}
}

// Write a type that implements StringEncoder
func (b StreamWriter) WriteEnc(v Encoder) {
	v.Encode(b)
}
