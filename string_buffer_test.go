package fast

import "testing"

func BenchmarkUint(b *testing.B) {
	var buf StringBuffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.WriteUint(123)
		buf.Reset()
	}
}

func BenchmarkUint64(b *testing.B) {
	var buf StringBuffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.WriteUint64(123)
		buf.Reset()
	}
}
