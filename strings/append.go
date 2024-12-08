package strings

import (
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"

	"github.com/webmafia/fast"
)

var digits [1000]uint32
var pow10 = [...]uint64{1, 10, 100, 1000, 10000, 100000, 1000000}

func init() {
	for i := uint32(0); i < 1000; i++ {
		digits[i] = (((i / 100) + '0') << 16) + ((((i / 10) % 10) + '0') << 8) + i%10 + '0'
		if i < 10 {
			digits[i] += 2 << 24
		} else if i < 100 {
			digits[i] += 1 << 24
		}
	}
}

func appendFirstBuf(b []byte, v uint32) []byte {
	start := v >> 24
	if start == 0 {
		b = append(b, byte(v>>16), byte(v>>8))
	} else if start == 1 {
		b = append(b, byte(v>>8))
	}
	b = append(b, byte(v))
	return b
}

func appendBuf(buf []byte, v uint32) []byte {
	return append(buf, byte(v>>16), byte(v>>8), byte(v))
}

func AppendBytes(b []byte, bytes []byte) []byte {
	return append(b, bytes...)
}

func AppendByte(b []byte, val byte) []byte {
	return append(b, val)
}

func AppendRune(b []byte, r rune) []byte {
	return utf8.AppendRune(b, r)
}

func AppendString(b []byte, s string) []byte {
	return append(b, s...)
}

func AppendBool(b []byte, v bool) []byte {
	if v {
		b = append(b, "true"...)
	} else {
		b = append(b, "false"...)
	}

	return b
}

func AppendUint8(b []byte, val uint8) []byte {
	return append(b, val)
}

func AppendInt8(b []byte, nval int8) []byte {
	var val uint8
	if nval < 0 {
		val = uint8(-nval)
		b = append(b, '-')
	} else {
		val = uint8(nval)
	}
	return append(b, val)
}

func AppendUint16(b []byte, val uint16) []byte {
	q1 := val / 1000
	if q1 == 0 {
		b = appendFirstBuf(b, digits[val])
		return b
	}
	r1 := val - q1*1000
	b = appendFirstBuf(b, digits[q1])
	b = appendBuf(b, digits[r1])
	return b
}

func AppendInt16(b []byte, nval int16) []byte {
	var val uint16
	if nval < 0 {
		val = uint16(-nval)
		b = append(b, '-')
	} else {
		val = uint16(nval)
	}
	b = AppendUint16(b, val)
	return b
}

func AppendUint32(b []byte, val uint32) []byte {
	q1 := val / 1000
	if q1 == 0 {
		b = appendFirstBuf(b, digits[val])
		return b
	}
	r1 := val - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		b = appendFirstBuf(b, digits[q1])
		b = appendBuf(b, digits[r1])
		return b
	}
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		b = appendFirstBuf(b, digits[q2])
	} else {
		r3 := q2 - q3*1000
		b = append(b, byte(q3+'0'))
		b = appendBuf(b, digits[r3])
	}
	b = appendBuf(b, digits[r2])
	b = appendBuf(b, digits[r1])
	return b
}

func AppendInt32(b []byte, nval int32) []byte {
	var val uint32
	if nval < 0 {
		val = uint32(-nval)
		b = append(b, '-')
	} else {
		val = uint32(nval)
	}
	b = AppendUint32(b, val)
	return b
}

func AppendUint64(b []byte, val uint64) []byte {
	q1 := val / 1000
	if q1 == 0 {
		b = appendFirstBuf(b, digits[val])
		return b
	}
	r1 := val - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		b = appendFirstBuf(b, digits[q1])
		b = appendBuf(b, digits[r1])
		return b
	}
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		b = appendFirstBuf(b, digits[q2])
		b = appendBuf(b, digits[r2])
		b = appendBuf(b, digits[r1])
		return b
	}
	r3 := q2 - q3*1000
	q4 := q3 / 1000
	if q4 == 0 {
		b = appendFirstBuf(b, digits[q3])
		b = appendBuf(b, digits[r3])
		b = appendBuf(b, digits[r2])
		b = appendBuf(b, digits[r1])
		return b
	}
	r4 := q3 - q4*1000
	q5 := q4 / 1000
	if q5 == 0 {
		b = appendFirstBuf(b, digits[q4])
		b = appendBuf(b, digits[r4])
		b = appendBuf(b, digits[r3])
		b = appendBuf(b, digits[r2])
		b = appendBuf(b, digits[r1])
		return b
	}
	r5 := q4 - q5*1000
	q6 := q5 / 1000
	if q6 == 0 {
		b = appendFirstBuf(b, digits[q5])
	} else {
		b = appendFirstBuf(b, digits[q6])
		r6 := q5 - q6*1000
		b = appendBuf(b, digits[r6])
	}
	b = appendBuf(b, digits[r5])
	b = appendBuf(b, digits[r4])
	b = appendBuf(b, digits[r3])
	b = appendBuf(b, digits[r2])
	b = appendBuf(b, digits[r1])
	return b
}

func AppendInt64(b []byte, nval int64) []byte {
	var val uint64
	if nval < 0 {
		val = uint64(-nval)
		b = append(b, '-')
	} else {
		val = uint64(nval)
	}
	b = AppendUint64(b, val)
	return b
}

func AppendInt(b []byte, val int) []byte {
	b = AppendInt64(b, int64(val))
	return b
}

func AppendUint(b []byte, val uint) []byte {
	b = AppendUint64(b, uint64(val))
	return b
}

func AppendFloat32(b []byte, val float32) []byte {
	abs := math.Abs(float64(val))
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if float32(abs) < 1e-6 || float32(abs) >= 1e21 {
			fmt = 'e'
		}
	}
	b = strconv.AppendFloat(b, float64(val), fmt, -1, 32)
	return b
}

func AppendFloat32Lossy(b []byte, val float32) []byte {
	if val < 0 {
		b = AppendByte(b, '-')
		val = -val
	}
	if val > 0x4ffffff {
		b = AppendFloat32(b, val)
		return b
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(float64(val)*float64(exp) + 0.5)
	b = AppendUint64(b, lval/exp)
	fval := lval % exp
	if fval == 0 {
		return b
	}
	b = AppendByte(b, '.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		b = AppendByte(b, '0')
	}
	b = AppendUint64(b, fval)
	for b[len(b)-1] == '0' {
		b = b[:len(b)-1]
	}
	return b
}

func AppendFloat64(b []byte, val float64) []byte {
	abs := math.Abs(val)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			fmt = 'e'
		}
	}
	b = strconv.AppendFloat(b, float64(val), fmt, -1, 64)
	return b
}

func AppendFloat64Lossy(b []byte, val float64) []byte {
	if val < 0 {
		b = AppendByte(b, '-')
		val = -val
	}
	if val > 0x4ffffff {
		b = AppendFloat64(b, val)
		return b
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(val*float64(exp) + 0.5)
	b = AppendUint64(b, lval/exp)
	fval := lval % exp
	if fval == 0 {
		return b
	}
	b = AppendByte(b, '.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		b = AppendByte(b, '0')
	}
	b = AppendUint64(b, fval)
	for b[len(b)-1] == '0' {
		b = b[:len(b)-1]
	}
	return b
}

func AppendStringer(b []byte, v fmt.Stringer) []byte {
	return AppendString(b, v.String())
}

func AppendTextAppender(b []byte, v fast.TextAppender) []byte {
	b, _ = v.AppendText(b)
	return b
}
