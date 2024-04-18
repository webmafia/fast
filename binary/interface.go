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

	WriteString(s string) (int, error)
	WriteUint8(v uint8) error
	WriteInt8(v int8) error
	WriteUint16(v uint16) error
	WriteInt16(v int16) error
	WriteUint32(v uint32) error
	WriteInt32(v int32) error
	WriteUint64(v uint64) error
	WriteInt64(v int64) error
	WriteInt(v int) error
	WriteUint(v uint) error
	WriteFloat32(v float32) error
	WriteFloat64(v float64) error
	WriteVarint(v int64) error
	WriteUvarint(v uint64) error
	WriteBool(v bool) error
	WriteEnc(v Encoder) error
	WriteVal(val any) error
}

type Encoder interface {
	Encode(w Writer) error
}
