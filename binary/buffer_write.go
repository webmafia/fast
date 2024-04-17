package binary

import (
	"github.com/webmafia/fast"
)

var _ Writer = (*BufferWriter)(nil)

// A Buffer is used to efficiently build a string using Write methods.
// It minimizes memory copying. The zero value is ready to use.
// Do not copy a non-zero Buffer.
type BufferWriter struct {
	buf []byte
}

func NewBufferWriter(cap int) *BufferWriter {
	return &BufferWriter{
		buf: make([]byte, 0, cap),
	}
}

// String returns the accumulated string.
func (b *BufferWriter) String() string {
	return fast.BytesToString(b.buf)
}

// String returns the accumulated string as bytes.
func (b *BufferWriter) Bytes() []byte {
	return b.buf
}

// Len returns the number of accumulated bytes; b.Len() == len(b.String()).
func (b *BufferWriter) Len() int {
	return len(b.buf)
}

// Cap returns the capacity of the Buffer's underlying byte slice. It is the
// total space allocated for the string being built and includes any bytes
// already written.
func (b *BufferWriter) Cap() int {
	return cap(b.buf)
}

// Reset resets the Buffer to be empty.
func (b *BufferWriter) Reset() {
	b.buf = b.buf[:0]
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *BufferWriter) grow(n int) {
	buf := fast.MakeNoZero(2*cap(b.buf) + n)[:len(b.buf)]
	copy(buf, b.buf)
	b.buf = buf
}

// Grow grows b's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to b
// without another allocation. If n is negative, Grow panics.
func (b *BufferWriter) Grow(n int) {
	if n < 0 {
		panic("fast.Buffer.Grow: negative count")
	}
	if cap(b.buf)-len(b.buf) < n {
		b.grow(n)
	}
}

// Write appends the contents of p to b's buffer.
// Write always returns len(p), nil.
func (b *BufferWriter) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}

// WriteByte appends the byte c to b's buffer.
// The returned error is always nil.
func (b *BufferWriter) WriteByte(c byte) error {
	b.buf = append(b.buf, c)
	return nil
}

// WriteString appends the contents of s to b's buffer.
// It returns the length of s and a nil error.
func (b *BufferWriter) WriteString(s string) (int, error) {
	b.buf = append(b.buf, s...)
	return len(s), nil
}

// Write a type that implements StringEncoder
func (b *BufferWriter) WriteEnc(v Encoder) {
	v.Encode(b)
}

func (b *BufferWriter) WriteVal(val any) {
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
