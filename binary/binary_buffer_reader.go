package fast

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/webmafia/fast"
)

var _ io.Reader = (*BinaryBufferReader)(nil)

type BinaryBufferReader struct {
	b      *BinaryBuffer
	cursor int
}

func NewBinaryBufferReader(b *BinaryBuffer) BinaryBufferReader {
	return BinaryBufferReader{
		b: b,
	}
}

func (b *BinaryBufferReader) Len() int {
	return b.b.Len()
}

func (b *BinaryBufferReader) ToRead() int {
	return b.Len() - b.cursor
}

func (b *BinaryBufferReader) Reset() {
	b.cursor = 0
}

func (b *BinaryBufferReader) Read(dst []byte) (n int, err error) {
	n = copy(dst, b.b.buf[b.cursor:])
	b.cursor += n
	return
}

func (b *BinaryBufferReader) ReadBytes(n int) (dst []byte) {
	dst = b.b.buf[b.cursor : b.cursor+n]
	b.cursor += n
	return
}

func (b *BinaryBufferReader) ReadString(n int) string {
	return fast.BytesToString(b.ReadBytes(n))
}

func (b *BinaryBufferReader) ReadBool() bool {
	v := b.b.buf[b.cursor]
	b.cursor++
	return v != 0
}

func (b *BinaryBufferReader) ReadUint8() (v uint8) {
	v = b.b.buf[b.cursor]
	b.cursor++
	return
}

func (b *BinaryBufferReader) ReadInt8() (v int8) {
	return int8(b.ReadUint8())
}

func (b *BinaryBufferReader) ReadUint16() (v uint16) {
	v = binary.LittleEndian.Uint16(b.b.buf[b.cursor:])
	b.cursor += 2
	return
}

// Write int16 (little endian)
func (b *BinaryBufferReader) ReadInt16() (v int16) {
	return int16(b.ReadUint16())
}

// Write uint32
func (b *BinaryBufferReader) ReadUint32() (v uint32) {
	v = binary.LittleEndian.Uint32(b.b.buf[b.cursor:])
	b.cursor += 4
	return
}

// Write int32
func (b *BinaryBufferReader) ReadInt32() (v int32) {
	return int32(b.ReadUint32())
}

// Write uint64
func (b *BinaryBufferReader) ReadUint64() (v uint64) {
	v = binary.LittleEndian.Uint64(b.b.buf[b.cursor:])
	b.cursor += 8
	return
}

// Write int64
func (b *BinaryBufferReader) ReadInt64() (v int64) {
	return int64(b.ReadUint64())
}

// Write int
func (b *BinaryBufferReader) ReadInt() (v int) {
	return int(b.ReadUint64())
}

// Write uint
func (b *BinaryBufferReader) ReadUint() (v uint) {
	return uint(b.ReadUint64())
}

// Write float32
func (b *BinaryBufferReader) ReadFloat32() (v float32) {
	return math.Float32frombits(b.ReadUint32())
}

// Write float64
func (b *BinaryBufferReader) ReadFloat64() (v float64) {
	return math.Float64frombits(b.ReadUint64())
}

func (b *BinaryBufferReader) ReadVarint() (v int64) {
	v, n := binary.Varint(b.b.buf[b.cursor:])
	b.cursor += n
	return
}

func (b *BinaryBufferReader) ReadUvarint() (v uint64) {
	v, n := binary.Uvarint(b.b.buf[b.cursor:])
	b.cursor += n
	return
}
