package fast

import (
	"fmt"
	"io"
)

var (
	_ io.Writer       = (*StringWriter)(nil)
	_ io.ByteWriter   = (*StringWriter)(nil)
	_ io.StringWriter = (*StringWriter)(nil)
	_ io.WriteCloser  = (*StringWriter)(nil)
	_ io.WriterTo     = (*StringWriter)(nil)
)

type StringWriter struct {
	w   io.Writer
	buf *StringBuffer
}

func NewStringWriter(w io.Writer, buf ...*StringBuffer) (sw StringWriter) {
	sw.w = w

	if len(buf) > 0 && buf[0] != nil {
		sw.buf = buf[0]
	} else {
		sw.buf = NewStringBuffer(4096)
	}

	return
}

func (w StringWriter) Available() int {
	return w.buf.Cap() - w.buf.Len()
}

func (w StringWriter) AvailableBuffer() []byte {
	return w.buf.buf[len(w.buf.buf):][:0]
}

func (w StringWriter) Flush() (err error) {
	if w.buf.Len() > 0 {
		_, err = w.w.Write(w.buf.Bytes())
		w.buf.Reset()
	}

	return
}

// Close implements io.WriteCloser.
func (w StringWriter) Close() (err error) {
	return w.Flush()
}

// WriteTo implements io.WriterTo.
func (w StringWriter) WriteTo(w2 io.Writer) (int64, error) {
	n, err := w2.Write(w.buf.buf)
	return int64(n), err
}

func (w StringWriter) maybeFlush() {
	if w.buf.Len() >= w.buf.Cap()/2 {
		w.Flush()
	}
}

func (w StringWriter) Write(p []byte) (n int, err error) {
	if len(p) > w.Available() {
		if err = w.Flush(); err != nil {
			return
		}

		if len(p) > w.buf.Cap() {
			return w.w.Write(p)
		}
	}

	return w.buf.Write(p)
}

func (w StringWriter) WriteByte(c byte) error {
	w.maybeFlush()
	return w.buf.WriteByte(c)
}

func (w StringWriter) WriteRune(r rune) (int, error) {
	w.maybeFlush()
	return w.buf.WriteRune(r)
}

func (w StringWriter) WriteString(s string) (int, error) {
	return w.Write(StringToBytes(s))
}

func (w StringWriter) WriteBool(v bool) {
	w.maybeFlush()
	w.buf.WriteBool(v)
}

func (w StringWriter) WriteUint8(val uint8) {
	w.maybeFlush()
	w.buf.WriteUint8(val)
}

func (w StringWriter) WriteInt8(val int8) {
	w.maybeFlush()
	w.buf.WriteInt8(val)
}

func (w StringWriter) WriteUint16(val uint16) {
	w.maybeFlush()
	w.buf.WriteUint16(val)
}

func (w StringWriter) WriteInt16(val int16) {
	w.maybeFlush()
	w.buf.WriteInt16(val)
}

func (w StringWriter) WriteUint32(val uint32) {
	w.maybeFlush()
	w.buf.WriteUint32(val)
}

func (w StringWriter) WriteInt32(val int32) {
	w.maybeFlush()
	w.buf.WriteInt32(val)
}

func (w StringWriter) WriteUint64(val uint64) {
	w.maybeFlush()
	w.buf.WriteUint64(val)
}

func (w StringWriter) WriteInt64(val int64) {
	w.maybeFlush()
	w.buf.WriteInt64(val)
}

func (w StringWriter) WriteInt(val int) {
	w.WriteInt64(int64(val))
}

func (w StringWriter) WriteUint(val uint) {
	w.WriteUint64(uint64(val))
}

func (w StringWriter) WriteFloat32(val float32) {
	w.maybeFlush()
	w.buf.WriteFloat32(val)
}

func (w StringWriter) WriteFloat32Lossy(val float32) {
	w.maybeFlush()
	w.buf.WriteFloat32Lossy(val)
}

func (w StringWriter) WriteFloat64(val float64) {
	w.maybeFlush()
	w.buf.WriteFloat64(val)
}

func (w StringWriter) WriteFloat64Lossy(val float64) {
	w.maybeFlush()
	w.buf.WriteFloat64Lossy(val)
}

func (w StringWriter) WriteVal(val any) (err error) {
	switch v := val.(type) {

	case TextAppender:
		var buf []byte

		if buf, err = v.AppendText(w.AvailableBuffer()); err != nil {
			return
		}

		_, err = w.Write(buf)

	case fmt.Stringer:
		_, err = w.WriteString(v.String())

	case string:
		_, err = w.WriteString(v)

	case []byte:
		_, err = w.Write(v)

	case int:
		w.WriteInt(v)

	case int8:
		w.WriteInt8(v)

	case int16:
		w.WriteInt16(v)

	case int32:
		w.WriteInt32(v)

	case int64:
		w.WriteInt64(v)

	case uint:
		w.WriteUint(v)

	case uint8:
		w.WriteUint8(v)

	case uint16:
		w.WriteUint16(v)

	case uint32:
		w.WriteUint32(v)

	case uint64:
		w.WriteUint64(v)

	case float32:
		w.WriteFloat32(v)

	case float64:
		w.WriteFloat64(v)

	case bool:
		w.WriteBool(v)

	default:
		return fmt.Errorf("unknown type: %T", v)

	}

	return
}
