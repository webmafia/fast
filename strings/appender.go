package strings

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/webmafia/fast"
)

var (
	ifaceTextAppender = reflect.TypeFor[fast.TextAppender]()
	ifaceStringer     = reflect.TypeFor[fmt.Stringer]()
)

type Appender func(b []byte, v unsafe.Pointer) []byte

func GetAppender(typ reflect.Type) (a Appender, err error) {
	if typ.Implements(ifaceTextAppender) {
		return AppendUnsafeTextAppender, nil
	}

	if typ.Implements(ifaceStringer) {
		return AppendUnsafeStringer, nil
	}

	switch kind := typ.Kind(); kind {

	case reflect.Bool:
		return AppendUnsafeBool, nil

	case reflect.Int:
		return AppendUnsafeInt, nil

	case reflect.Int8:
		return AppendUnsafeInt8, nil

	case reflect.Int16:
		return AppendUnsafeInt16, nil

	case reflect.Int32:
		return AppendUnsafeInt32, nil

	case reflect.Int64:
		return AppendUnsafeInt64, nil

	case reflect.Uint:
		return AppendUnsafeUint, nil

	case reflect.Uint8:
		return AppendUnsafeUint8, nil

	case reflect.Uint16:
		return AppendUnsafeUint16, nil

	case reflect.Uint32:
		return AppendUnsafeUint32, nil

	case reflect.Uint64:
		return AppendUnsafeUint64, nil

	case reflect.Float32:
		return AppendUnsafeFloat32Lossy, nil

	case reflect.Float64:
		return AppendUnsafeFloat64Lossy, nil

	case reflect.String:
		return AppendUnsafeString, nil

	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			return AppendUnsafeBytes, nil
		}

		fallthrough

	default:
		return nil, fmt.Errorf("unsupported type: %s", kind)

	}
}
