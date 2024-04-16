// Borrowed from jsoniter (https://github.com/json-iterator/go)
package fast

import (
	"encoding/binary"
	"math"
)

// Write uint8
func (b BinaryStream) WriteUint8(v uint8) {
	b.WriteByte(v)
}

// Write int8
func (b BinaryStream) WriteInt8(v int8) {
	b.WriteByte(uint8(v))
}

// Write uint16
func (b BinaryStream) WriteUint16(v uint16) {
	buf := b.buf.AvailableBuffer()
	buf = binary.LittleEndian.AppendUint16(buf, v)
	b.buf.Write(buf)
}

// Write int16
func (b BinaryStream) WriteInt16(v int16) {
	b.WriteUint16(uint16(v))
}

// Write uint32
func (b BinaryStream) WriteUint32(v uint32) {
	buf := b.buf.AvailableBuffer()
	buf = binary.LittleEndian.AppendUint32(buf, v)
	b.buf.Write(buf)
}

// Write int32
func (b BinaryStream) WriteInt32(v int32) {
	b.WriteUint32(uint32(v))
}

// Write uint64
func (b BinaryStream) WriteUint64(v uint64) {
	buf := b.buf.AvailableBuffer()
	buf = binary.LittleEndian.AppendUint64(buf, v)
	b.buf.Write(buf)
}

// Write int64
func (b BinaryStream) WriteInt64(v int64) {
	b.WriteUint64(uint64(v))
}

// Write int
func (b BinaryStream) WriteInt(v int) {
	b.WriteInt64(int64(v))
}

// Write uint
func (b BinaryStream) WriteUint(v uint) {
	b.WriteUint64(uint64(v))
}

// Write float32
func (b BinaryStream) WriteFloat32(v float32) {
	b.WriteUint32(math.Float32bits(v))
}

// Write float64
func (b BinaryStream) WriteFloat64(v float64) {
	b.WriteUint64(math.Float64bits(v))
}

func (b BinaryStream) WriteVarint(v int64) {
	buf := b.buf.AvailableBuffer()
	buf = binary.AppendVarint(buf, v)
	b.buf.Write(buf)
}

func (b BinaryStream) WriteUvarint(v uint64) {
	buf := b.buf.AvailableBuffer()
	buf = binary.AppendUvarint(buf, v)
	b.buf.Write(buf)
}
