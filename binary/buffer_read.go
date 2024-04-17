package binary

import (
	"github.com/webmafia/fast"
)

var _ Reader = (*BufferReader)(nil)

type BufferReader struct {
	buf    []byte
	cursor int
}

func NewBufferReader(buf []byte) BufferReader {
	return BufferReader{
		buf: buf,
	}
}

func (b *BufferReader) Len() int {
	return b.Cap() - b.cursor
}

func (b *BufferReader) Cap() int {
	return len(b.buf)
}

func (b *BufferReader) Reset() {
	b.cursor = 0
}

func (b *BufferReader) Read(dst []byte) (n int, err error) {
	n = copy(dst, b.buf[b.cursor:])
	b.cursor += n
	return
}

func (b *BufferReader) ReadByte() (byte, error) {
	v := b.buf[b.cursor]
	b.cursor++
	return v, nil
}

func (b *BufferReader) ReadBytes(n int) (dst []byte) {
	dst = b.buf[b.cursor : b.cursor+n]
	b.cursor += n
	return
}

func (b *BufferReader) ReadString(n int) string {
	return fast.BytesToString(b.ReadBytes(n))
}
