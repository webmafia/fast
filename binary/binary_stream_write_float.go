// Borrowed from jsoniter (https://github.com/json-iterator/go)
package fast

import (
	"math"
)

// Write float32
func (b BinaryStreamWriter) WriteFloat32(v float32) {
	b.WriteUint32(math.Float32bits(v))
}

// Write float64
func (b BinaryStreamWriter) WriteFloat64(v float64) {
	b.WriteUint64(math.Float64bits(v))
}
