package fast

import (
	"testing"
)

func BenchmarkVarintWrite(b *testing.B) {
	buf := NewBinaryBuffer(10)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.WriteVarint(int64(b.N))
		buf.Reset()
	}
}

func BenchmarkVarintRead(b *testing.B) {
	buf := NewBinaryBuffer(10)
	buf.WriteVarint(int64(b.N))
	r := NewBinaryBufferReader(buf)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = r.ReadVarint()
		r.Reset()
	}
}
