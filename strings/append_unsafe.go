package strings

import (
	"fmt"
	"unsafe"

	"github.com/webmafia/fast"
)

func AppendUnsafeBytes(b []byte, v unsafe.Pointer) []byte {
	return AppendBytes(b, *(*[]byte)(v))
}

func AppendUnsafeByte(b []byte, v unsafe.Pointer) []byte {
	return AppendByte(b, *(*byte)(v))
}

func AppendUnsafeRune(b []byte, v unsafe.Pointer) []byte {
	return AppendRune(b, *(*rune)(v))
}

func AppendUnsafeString(b []byte, v unsafe.Pointer) []byte {
	return AppendString(b, *(*string)(v))
}

func AppendUnsafeBool(b []byte, v unsafe.Pointer) []byte {
	return AppendBool(b, *(*bool)(v))
}

func AppendUnsafeUint8(b []byte, v unsafe.Pointer) []byte {
	return AppendUint8(b, *(*uint8)(v))
}

func AppendUnsafeInt8(b []byte, v unsafe.Pointer) []byte {
	return AppendInt8(b, *(*int8)(v))
}

func AppendUnsafeUint16(b []byte, v unsafe.Pointer) []byte {
	return AppendUint16(b, *(*uint16)(v))
}

func AppendUnsafeInt16(b []byte, v unsafe.Pointer) []byte {
	return AppendInt16(b, *(*int16)(v))
}

func AppendUnsafeUint32(b []byte, v unsafe.Pointer) []byte {
	return AppendUint32(b, *(*uint32)(v))
}

func AppendUnsafeInt32(b []byte, v unsafe.Pointer) []byte {
	return AppendInt32(b, *(*int32)(v))
}

func AppendUnsafeUint64(b []byte, v unsafe.Pointer) []byte {
	return AppendUint64(b, *(*uint64)(v))
}

func AppendUnsafeInt64(b []byte, v unsafe.Pointer) []byte {
	return AppendInt64(b, *(*int64)(v))
}

func AppendUnsafeInt(b []byte, v unsafe.Pointer) []byte {
	return AppendInt(b, *(*int)(v))
}

func AppendUnsafeUint(b []byte, v unsafe.Pointer) []byte {
	return AppendUint(b, *(*uint)(v))
}

func AppendUnsafeFloat32(b []byte, v unsafe.Pointer) []byte {
	return AppendFloat32(b, *(*float32)(v))
}

func AppendUnsafeFloat32Lossy(b []byte, v unsafe.Pointer) []byte {
	return AppendFloat32Lossy(b, *(*float32)(v))
}

func AppendUnsafeFloat64(b []byte, v unsafe.Pointer) []byte {
	return AppendFloat64(b, *(*float64)(v))
}

func AppendUnsafeFloat64Lossy(b []byte, v unsafe.Pointer) []byte {
	return AppendFloat64Lossy(b, *(*float64)(v))
}

func AppendUnsafeStringer(b []byte, v unsafe.Pointer) []byte {
	return AppendStringer(b, *(*fmt.Stringer)(v))
}

func AppendUnsafeTextAppender(b []byte, v unsafe.Pointer) []byte {
	return AppendTextAppender(b, *(*fast.TextAppender)(v))
}
