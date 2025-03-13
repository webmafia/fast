package buffer

import (
	"math"
	"strconv"
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

func writeFirstBuf(space []byte, v uint32) []byte {
	start := v >> 24
	if start == 0 {
		space = append(space, byte(v>>16), byte(v>>8))
	} else if start == 1 {
		space = append(space, byte(v>>8))
	}
	space = append(space, byte(v))
	return space
}

func writeBuf(buf []byte, v uint32) []byte {
	return append(buf, byte(v>>16), byte(v>>8), byte(v))
}

// Write uint8
func (b StringBuffer) WriteUint8(val uint8) {
	b.B.B = writeFirstBuf(b.B.B, digits[val])
}

// Write int8
func (b StringBuffer) WriteInt8(nval int8) {
	var val uint8
	if nval < 0 {
		val = uint8(-nval)
		b.B.B = append(b.B.B, '-')
	} else {
		val = uint8(nval)
	}
	b.B.B = writeFirstBuf(b.B.B, digits[val])
}

// Write uint16
func (b StringBuffer) WriteUint16(val uint16) {
	q1 := val / 1000
	if q1 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[val])
		return
	}
	r1 := val - q1*1000
	b.B.B = writeFirstBuf(b.B.B, digits[q1])
	b.B.B = writeBuf(b.B.B, digits[r1])
}

// Write int16
func (b StringBuffer) WriteInt16(nval int16) {
	var val uint16
	if nval < 0 {
		val = uint16(-nval)
		b.B.B = append(b.B.B, '-')
	} else {
		val = uint16(nval)
	}
	b.WriteUint16(val)
}

// Write uint32
func (b StringBuffer) WriteUint32(val uint32) {
	q1 := val / 1000
	if q1 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[val])
		return
	}
	r1 := val - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[q1])
		b.B.B = writeBuf(b.B.B, digits[r1])
		return
	}
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[q2])
	} else {
		r3 := q2 - q3*1000
		b.B.B = append(b.B.B, byte(q3+'0'))
		b.B.B = writeBuf(b.B.B, digits[r3])
	}
	b.B.B = writeBuf(b.B.B, digits[r2])
	b.B.B = writeBuf(b.B.B, digits[r1])
}

// Write int32
func (b StringBuffer) WriteInt32(nval int32) {
	var val uint32
	if nval < 0 {
		val = uint32(-nval)
		b.B.B = append(b.B.B, '-')
	} else {
		val = uint32(nval)
	}
	b.WriteUint32(val)
}

// Write uint64
func (b StringBuffer) WriteUint64(val uint64) {
	q1 := val / 1000
	if q1 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[val])
		return
	}
	r1 := val - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[q1])
		b.B.B = writeBuf(b.B.B, digits[r1])
		return
	}
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[q2])
		b.B.B = writeBuf(b.B.B, digits[r2])
		b.B.B = writeBuf(b.B.B, digits[r1])
		return
	}
	r3 := q2 - q3*1000
	q4 := q3 / 1000
	if q4 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[q3])
		b.B.B = writeBuf(b.B.B, digits[r3])
		b.B.B = writeBuf(b.B.B, digits[r2])
		b.B.B = writeBuf(b.B.B, digits[r1])
		return
	}
	r4 := q3 - q4*1000
	q5 := q4 / 1000
	if q5 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[q4])
		b.B.B = writeBuf(b.B.B, digits[r4])
		b.B.B = writeBuf(b.B.B, digits[r3])
		b.B.B = writeBuf(b.B.B, digits[r2])
		b.B.B = writeBuf(b.B.B, digits[r1])
		return
	}
	r5 := q4 - q5*1000
	q6 := q5 / 1000
	if q6 == 0 {
		b.B.B = writeFirstBuf(b.B.B, digits[q5])
	} else {
		b.B.B = writeFirstBuf(b.B.B, digits[q6])
		r6 := q5 - q6*1000
		b.B.B = writeBuf(b.B.B, digits[r6])
	}
	b.B.B = writeBuf(b.B.B, digits[r5])
	b.B.B = writeBuf(b.B.B, digits[r4])
	b.B.B = writeBuf(b.B.B, digits[r3])
	b.B.B = writeBuf(b.B.B, digits[r2])
	b.B.B = writeBuf(b.B.B, digits[r1])
}

// Write int64
func (b StringBuffer) WriteInt64(nval int64) {
	var val uint64
	if nval < 0 {
		val = uint64(-nval)
		b.B.B = append(b.B.B, '-')
	} else {
		val = uint64(nval)
	}
	b.WriteUint64(val)
}

// Write int
func (b StringBuffer) WriteInt(val int) {
	b.WriteInt64(int64(val))
}

// Write uint
func (b StringBuffer) WriteUint(val uint) {
	b.WriteUint64(uint64(val))
}

// Write float32
func (b StringBuffer) WriteFloat32(val float32) {
	abs := math.Abs(float64(val))
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if float32(abs) < 1e-6 || float32(abs) >= 1e21 {
			fmt = 'e'
		}
	}
	b.B.B = strconv.AppendFloat(b.B.B, float64(val), fmt, -1, 32)
}

// Write float32 with ONLY 6 digits precision although much much faster
func (b StringBuffer) WriteFloat32Lossy(val float32) {
	if val < 0 {
		b.B.WriteByte('-')
		val = -val
	}
	if val > 0x4ffffff {
		b.WriteFloat32(val)
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(float64(val)*float64(exp) + 0.5)
	b.WriteUint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	b.B.WriteByte('.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		b.B.WriteByte('0')
	}
	b.WriteUint64(fval)
	for b.B.B[len(b.B.B)-1] == '0' {
		b.B.B = b.B.B[:len(b.B.B)-1]
	}
}

// Write float64
func (b StringBuffer) WriteFloat64(val float64) {
	abs := math.Abs(val)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			fmt = 'e'
		}
	}
	b.B.B = strconv.AppendFloat(b.B.B, float64(val), fmt, -1, 64)
}

// Write float64 with ONLY 6 digits precision although much much faster
func (b StringBuffer) WriteFloat64Lossy(val float64) {
	if val < 0 {
		b.B.WriteByte('-')
		val = -val
	}
	if val > 0x4ffffff {
		b.WriteFloat64(val)
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(val*float64(exp) + 0.5)
	b.WriteUint64(lval / exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	b.B.WriteByte('.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		b.B.WriteByte('0')
	}
	b.WriteUint64(fval)
	for b.B.B[len(b.B.B)-1] == '0' {
		b.B.B = b.B.B[:len(b.B.B)-1]
	}
}
