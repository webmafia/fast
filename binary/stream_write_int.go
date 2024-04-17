// Borrowed from jsoniter (https://github.com/json-iterator/go)
package binary

import (
	"encoding/binary"
)

// Write int8
func (b StreamWriter) WriteInt8(v int8) {
	b.WriteByte(uint8(v))
}

// Write int16
func (b StreamWriter) WriteInt16(v int16) {
	b.WriteUint16(uint16(v))
}

// Write int32
func (b StreamWriter) WriteInt32(v int32) {
	b.WriteUint32(uint32(v))
}

// Write int64
func (b StreamWriter) WriteInt64(v int64) {
	b.WriteUint64(uint64(v))
}

// Write int
func (b StreamWriter) WriteInt(v int) {
	b.WriteInt64(int64(v))
}

func (b StreamWriter) WriteVarint(v int64) {
	buf := b.buf.AvailableBuffer()
	buf = binary.AppendVarint(buf, v)
	b.buf.Write(buf)
}
