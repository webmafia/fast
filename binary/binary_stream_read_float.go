package fast

import (
	"math"
)

func (b *BinaryStreamReader) ReadFloat32() float32 {
	return math.Float32frombits(b.ReadUint32())
}

func (b *BinaryStreamReader) ReadFloat64() float64 {
	return math.Float64frombits(b.ReadUint64())
}
