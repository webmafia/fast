package binary

import (
	"io"
	"testing"
)

func BenchmarkStreamWriter(b *testing.B) {
	w := NewStreamWriter(io.Discard)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.WriteInt(i)
	}
}
