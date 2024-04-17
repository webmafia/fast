// Borrowed from jsoniter (https://github.com/json-iterator/go)
package fast

import (
	"encoding/binary"
)

// Write int8
func (b BinaryStreamWriter) WriteInt8(v int8) {
	b.WriteByte(uint8(v))
}

// Write int16
func (b BinaryStreamWriter) WriteInt16(v int16) {
	b.WriteUint16(uint16(v))
}

// Write int32
func (b BinaryStreamWriter) WriteInt32(v int32) {
	b.WriteUint32(uint32(v))
}

// Write int64
func (b BinaryStreamWriter) WriteInt64(v int64) {
	b.WriteUint64(uint64(v))
}

// Write int
func (b BinaryStreamWriter) WriteInt(v int) {
	b.WriteInt64(int64(v))
}

func (b BinaryStreamWriter) WriteVarint(v int64) {
	buf := b.buf.AvailableBuffer()
	buf = binary.AppendVarint(buf, v)
	b.buf.Write(buf)
}
