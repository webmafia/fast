package binary

import (
	"encoding/binary"
)

func (b *BufferReader) ReadUint8() (v uint8) {
	v = b.buf[b.cursor]
	b.cursor++
	return
}

func (b *BufferReader) ReadUint16() (v uint16) {
	v = binary.LittleEndian.Uint16(b.buf[b.cursor:])
	b.cursor += 2
	return
}

// Write uint32
func (b *BufferReader) ReadUint32() (v uint32) {
	v = binary.LittleEndian.Uint32(b.buf[b.cursor:])
	b.cursor += 4
	return
}

// Write uint64
func (b *BufferReader) ReadUint64() (v uint64) {
	v = binary.LittleEndian.Uint64(b.buf[b.cursor:])
	b.cursor += 8
	return
}

// Write uint
func (b *BufferReader) ReadUint() (v uint) {
	return uint(b.ReadUint64())
}

func (b *BufferReader) ReadUvarint() (v uint64) {
	v, n := binary.Uvarint(b.buf[b.cursor:])
	b.cursor += n
	return
}
