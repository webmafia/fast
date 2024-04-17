package binary

import (
	"encoding/binary"
)

func (b *StreamReader) ReadInt8() int8 {
	var v uint8
	v, b.err = b.ReadByte()
	return int8(v)
}

func (b *StreamReader) ReadInt16() int16 {
	var v [2]byte
	b.err = b.ReadFull(v[:])
	return int16(binary.LittleEndian.Uint16(v[:]))
}

func (b *StreamReader) ReadInt32() int32 {
	var v [4]byte
	b.err = b.ReadFull(v[:])
	return int32(binary.LittleEndian.Uint32(v[:]))
}

func (b *StreamReader) ReadInt64() int64 {
	var v [8]byte
	b.err = b.ReadFull(v[:])
	return int64(binary.LittleEndian.Uint64(v[:]))
}

func (b *StreamReader) ReadInt() int {
	return int(b.ReadInt64())
}

func (b *StreamReader) ReadVarint() (v int64) {
	v, b.err = binary.ReadVarint(b.buf)
	return
}
