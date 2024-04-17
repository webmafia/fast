package binary

import (
	"math"
)

func (b *StreamReader) ReadFloat32() float32 {
	return math.Float32frombits(b.ReadUint32())
}

func (b *StreamReader) ReadFloat64() float64 {
	return math.Float64frombits(b.ReadUint64())
}
