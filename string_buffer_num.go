/*
	Borrowed from jsoniter (https://github.com/json-iterator/go)
*/

package fast

import "strconv"

var digits [1000]uint32

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

// WriteUint64 write uint64 to stream
func (stream *StringBuffer) WriteUint64(val uint64) {
	q1 := val / 1000
	if q1 == 0 {
		stream.buf = writeFirstBuf(stream.buf, digits[val])
		return
	}
	r1 := val - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		stream.buf = writeFirstBuf(stream.buf, digits[q1])
		stream.buf = writeBuf(stream.buf, digits[r1])
		return
	}
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		stream.buf = writeFirstBuf(stream.buf, digits[q2])
		stream.buf = writeBuf(stream.buf, digits[r2])
		stream.buf = writeBuf(stream.buf, digits[r1])
		return
	}
	r3 := q2 - q3*1000
	q4 := q3 / 1000
	if q4 == 0 {
		stream.buf = writeFirstBuf(stream.buf, digits[q3])
		stream.buf = writeBuf(stream.buf, digits[r3])
		stream.buf = writeBuf(stream.buf, digits[r2])
		stream.buf = writeBuf(stream.buf, digits[r1])
		return
	}
	r4 := q3 - q4*1000
	q5 := q4 / 1000
	if q5 == 0 {
		stream.buf = writeFirstBuf(stream.buf, digits[q4])
		stream.buf = writeBuf(stream.buf, digits[r4])
		stream.buf = writeBuf(stream.buf, digits[r3])
		stream.buf = writeBuf(stream.buf, digits[r2])
		stream.buf = writeBuf(stream.buf, digits[r1])
		return
	}
	r5 := q4 - q5*1000
	q6 := q5 / 1000
	if q6 == 0 {
		stream.buf = writeFirstBuf(stream.buf, digits[q5])
	} else {
		stream.buf = writeFirstBuf(stream.buf, digits[q6])
		r6 := q5 - q6*1000
		stream.buf = writeBuf(stream.buf, digits[r6])
	}
	stream.buf = writeBuf(stream.buf, digits[r5])
	stream.buf = writeBuf(stream.buf, digits[r4])
	stream.buf = writeBuf(stream.buf, digits[r3])
	stream.buf = writeBuf(stream.buf, digits[r2])
	stream.buf = writeBuf(stream.buf, digits[r1])
}

// TODO: Replace functions below with functions from jsoniter

func (b *StringBuffer) WriteInt(i int64) {
	b.buf = strconv.AppendInt(b.buf, i, 10)
}

func (b *StringBuffer) WriteUint(i uint64) {
	b.buf = strconv.AppendUint(b.buf, i, 10)
}

func (b *StringBuffer) WriteFloat64(f float64, dec int) {
	b.buf = strconv.AppendFloat(b.buf, f, 'f', dec, 64)
}

func (b *StringBuffer) WriteFloat32(f float32, dec int) {
	b.buf = strconv.AppendFloat(b.buf, float64(f), 'f', dec, 32)
}
