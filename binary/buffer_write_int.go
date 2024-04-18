// Borrowed from jsoniter (https://github.com/json-iterator/go)
package binary

import (
	"encoding/binary"
)

// Write int8
func (b *BufferWriter) WriteInt8(v int8) error {
	b.WriteByte(uint8(v))
	return nil
}

// Write int16
func (b *BufferWriter) WriteInt16(v int16) error {
	b.WriteUint16(uint16(v))
	return nil
}

// Write int32
func (b *BufferWriter) WriteInt32(v int32) error {
	b.WriteUint32(uint32(v))
	return nil
}

// Write int64
func (b *BufferWriter) WriteInt64(v int64) error {
	b.WriteUint64(uint64(v))
	return nil
}

// Write int
func (b *BufferWriter) WriteInt(v int) error {
	b.WriteInt64(int64(v))
	return nil
}

func (b *BufferWriter) WriteVarint(v int64) error {
	b.buf = binary.AppendVarint(b.buf, v)
	return nil
}
