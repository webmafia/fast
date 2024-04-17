package binary

import (
	"encoding/binary"
)

func (b *BufferReader) ReadInt8() (v int8) {
	return int8(b.ReadUint8())
}

// Write int16 (little endian)
func (b *BufferReader) ReadInt16() (v int16) {
	return int16(b.ReadUint16())
}

// Write int32
func (b *BufferReader) ReadInt32() (v int32) {
	return int32(b.ReadUint32())
}

// Write int64
func (b *BufferReader) ReadInt64() (v int64) {
	return int64(b.ReadUint64())
}

// Write int
func (b *BufferReader) ReadInt() (v int) {
	return int(b.ReadUint64())
}

func (b *BufferReader) ReadVarint() (v int64) {
	v, n := binary.Varint(b.buf[b.cursor:])
	b.cursor += n
	return
}
