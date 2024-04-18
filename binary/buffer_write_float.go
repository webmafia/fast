// Borrowed from jsoniter (https://github.com/json-iterator/go)
package binary

import (
	"math"
)

// Write float32
func (b *BufferWriter) WriteFloat32(v float32) error {
	return b.WriteUint32(math.Float32bits(v))
}

// Write float64
func (b *BufferWriter) WriteFloat64(v float64) error {
	return b.WriteUint64(math.Float64bits(v))
}
