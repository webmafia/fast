// Borrowed from jsoniter (https://github.com/json-iterator/go)
package binary

import (
	"encoding/binary"
)

// Write int8
func (b *StreamWriter) WriteInt8(v int8) error {
	return b.WriteByte(uint8(v))
}

// Write int16
func (b *StreamWriter) WriteInt16(v int16) error {
	return b.WriteUint16(uint16(v))
}

// Write int32
func (b *StreamWriter) WriteInt32(v int32) error {
	return b.WriteUint32(uint32(v))
}

// Write int64
func (b *StreamWriter) WriteInt64(v int64) error {
	return b.WriteUint64(uint64(v))
}

// Write int
func (b *StreamWriter) WriteInt(v int) error {
	return b.WriteInt64(int64(v))
}

func (b *StreamWriter) WriteVarint(v int64) (err error) {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], v)
	_, err = b.Write(buf[:n])
	return
}
