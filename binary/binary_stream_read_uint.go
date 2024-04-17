package fast

import (
	"encoding/binary"
)

func (b *BinaryStreamReader) ReadUint8() (v uint8) {
	v, b.err = b.ReadByte()
	return
}

func (b *BinaryStreamReader) ReadUint16() uint16 {
	var v [2]byte
	_, b.err = b.Read(v[:])
	return binary.LittleEndian.Uint16(v[:])
}

func (b *BinaryStreamReader) ReadUint32() uint32 {
	var v [4]byte
	_, b.err = b.Read(v[:])
	return binary.LittleEndian.Uint32(v[:])
}

func (b *BinaryStreamReader) ReadUint64() uint64 {
	var v [8]byte
	_, b.err = b.Read(v[:])
	return binary.LittleEndian.Uint64(v[:])
}

func (b *BinaryStreamReader) ReadUint() uint {
	return uint(b.ReadUint64())
}

func (b *BinaryStreamReader) ReadUvarint() (v uint64) {
	v, b.err = binary.ReadUvarint(b.buf)
	return
}
