package binary

import (
	"math"
)

// Write float32
func (b *BufferReader) ReadFloat32() (v float32) {
	return math.Float32frombits(b.ReadUint32())
}

// Write float64
func (b *BufferReader) ReadFloat64() (v float64) {
	return math.Float64frombits(b.ReadUint64())
}
