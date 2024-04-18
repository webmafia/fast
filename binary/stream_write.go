package binary

import (
	"io"

	"github.com/webmafia/fast"
)

var _ Writer = (*StreamWriter)(nil)

// A Stream is used to efficiently write binary data to
// an io.Writer.
type StreamWriter struct {
	w io.Writer
}

func NewStreamWriter(w io.Writer) *StreamWriter {
	return &StreamWriter{
		w: w,
	}
}

// Reset resets the Stream to be empty.
func (b *StreamWriter) Reset(w io.Writer) {
	b.w = w
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
//
//go:inline
func (b *StreamWriter) Write(p []byte) (int, error) {
	return b.w.Write(fast.NoescapeBytes(p))
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *StreamWriter) WriteByte(c byte) (err error) {
	_, err = b.Write([]byte{c})
	return
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *StreamWriter) WriteString(s string) (int, error) {
	return b.Write(fast.StringToBytes(s))
}

func (b *StreamWriter) WriteVal(val any) error {
	switch v := val.(type) {

	case Encoder:
		return b.WriteEnc(v)

	case string:
		_, err := b.WriteString(v)
		return err

	case []byte:
		_, err := b.Write(v)
		return err

	case int:
		return b.WriteInt(v)

	case int8:
		return b.WriteInt8(v)

	case int16:
		return b.WriteInt16(v)

	case int32:
		return b.WriteInt32(v)

	case int64:
		return b.WriteInt64(v)

	case uint:
		return b.WriteUint(v)

	case uint8:
		return b.WriteUint8(v)

	case uint16:
		return b.WriteUint16(v)

	case uint32:
		return b.WriteUint32(v)

	case uint64:
		return b.WriteUint64(v)

	case float32:
		return b.WriteFloat32(v)

	case float64:
		return b.WriteFloat64(v)

	case bool:
		return b.WriteBool(v)

	}

	return ErrUnknownValue
}

// Write a type that implements StringEncoder
func (b *StreamWriter) WriteEnc(v Encoder) error {
	return v.Encode(b)
}
