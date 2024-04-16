package fast

import "io"

type Writer interface {
	io.Writer
	io.ByteWriter

	Len() int
	Cap() int
	WriteRune(r rune) (int, error)
	WriteString(s string) (int, error)
	WriteUint8(v uint8)
	WriteInt8(v int8)
	WriteUint16(v uint16)
	WriteInt16(v int16)
	WriteUint32(v uint32)
	WriteInt32(v int32)
	WriteUint64(v uint64)
	WriteInt64(v int64)
	WriteInt(v int)
	WriteUint(v uint)
	WriteFloat32(v float32)
	WriteFloat64(v float64)
	WriteVarint(v int64)
	WriteUvarint(v uint64)
	WriteBool(v bool)
	WriteEnc(v BinaryEncoder)
	WriteVal(val any)
}

type BinaryEncoder interface {
	BinaryEncode(w Writer)
}
