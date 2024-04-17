package binary

import "io"

type Reader interface {
	io.Reader
	io.ByteReader

	Len() int
	Cap() int
	ReadBytes(n int) []byte
	ReadString(n int) string
	ReadUint8() uint8
	ReadInt8() int8
	ReadUint16() uint16
	ReadInt16() int16
	ReadUint32() uint32
	ReadInt32() int32
	ReadUint64() uint64
	ReadInt64() int64
	ReadInt() int
	ReadUint() uint
	ReadFloat32() float32
	ReadFloat64() float64
	ReadVarint() int64
	ReadUvarint() uint64
	ReadBool() bool
}

type Writer interface {
	io.Writer
	io.ByteWriter

	Len() int
	Cap() int
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
	WriteEnc(v Encoder)
	WriteVal(val any)
}

type Encoder interface {
	Encode(w Writer)
}
