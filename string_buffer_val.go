package fast

import "fmt"

func (b *StringBuffer) WriteVal(val any) {
	switch v := val.(type) {

	case StringEncoder:
		b.WriteEnc(v)

	case fmt.Stringer:
		b.WriteString(v.String())

	case string:
		b.WriteString(v)

	case []byte:
		b.Write(v)

	case int:
		b.WriteInt(v)

	case int8:
		b.WriteInt8(v)

	case int16:
		b.WriteInt16(v)

	case int32:
		b.WriteInt32(v)

	case int64:
		b.WriteInt64(v)

	case uint:
		b.WriteUint(v)

	case uint8:
		b.WriteUint8(v)

	case uint16:
		b.WriteUint16(v)

	case uint32:
		b.WriteUint32(v)

	case uint64:
		b.WriteUint64(v)

	case float32:
		b.WriteFloat32(v)

	case float64:
		b.WriteFloat64(v)

	case bool:
		b.WriteBool(v)

	default:
		b.WriteString("(any)")

	}
}
