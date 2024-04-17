package binary

import (
	"encoding/binary"
)

func (b *StreamReader) ReadUint8() (v uint8) {
	v, b.err = b.ReadByte()
	return
}

func (b *StreamReader) ReadUint16() uint16 {
	var v [2]byte
	b.err = b.ReadFull(v[:])
	return binary.LittleEndian.Uint16(v[:])
}

func (b *StreamReader) ReadUint32() uint32 {
	var v [4]byte
	b.err = b.ReadFull(v[:])
	return binary.LittleEndian.Uint32(v[:])
}

func (b *StreamReader) ReadUint64() uint64 {
	var v [8]byte
	b.err = b.ReadFull(v[:])
	return binary.LittleEndian.Uint64(v[:])
}

func (b *StreamReader) ReadUint() uint {
	return uint(b.ReadUint64())
}

func (b *StreamReader) ReadUvarint() (v uint64) {
	v, b.err = binary.ReadUvarint(b.buf)
	return
}
