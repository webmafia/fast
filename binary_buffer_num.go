// Borrowed from jsoniter (https://github.com/json-iterator/go)
package fast

import (
	"encoding/binary"
	"math"
)

// Write uint8
func (b *BinaryBuffer) WriteUint8(v uint8) {
	b.WriteByte(v)
}

// Write int8
func (b *BinaryBuffer) WriteInt8(v int8) {
	b.WriteByte(uint8(v))
}

// Write uint16 (little endian)
func (b *BinaryBuffer) WriteUint16(v uint16) {
	b.buf = append(b.buf,
		byte(v),
		byte(v>>8),
	)
}

// Write int16 (little endian)
func (b *BinaryBuffer) WriteInt16(v int16) {
	b.buf = append(b.buf,
		byte(v),
		byte(v>>8),
	)
}

// Write uint32
func (b *BinaryBuffer) WriteUint32(v uint32) {
	b.buf = append(b.buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

// Write int32
func (b *BinaryBuffer) WriteInt32(v int32) {
	b.buf = append(b.buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
	)
}

// Write uint64
func (b *BinaryBuffer) WriteUint64(v uint64) {
	b.buf = append(b.buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
}

// Write int64
func (b *BinaryBuffer) WriteInt64(v int64) {
	b.buf = append(b.buf,
		byte(v),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56),
	)
}

// Write int
func (b *BinaryBuffer) WriteInt(v int) {
	b.WriteInt64(int64(v))
}

// Write uint
func (b *BinaryBuffer) WriteUint(v uint) {
	b.WriteUint64(uint64(v))
}

// Write float32
func (b *BinaryBuffer) WriteFloat32(v float32) {
	b.WriteUint32(math.Float32bits(v))
}

// Write float64
func (b *BinaryBuffer) WriteFloat64(v float64) {
	b.WriteUint64(math.Float64bits(v))
}

func (b *BinaryBuffer) WriteVarint(v int64) {
	b.buf = binary.AppendVarint(b.buf, v)
}

func (b *BinaryBuffer) WriteUvarint(v uint64) {
	b.buf = binary.AppendUvarint(b.buf, v)
}
