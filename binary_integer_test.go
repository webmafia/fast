package fast

import (
	"encoding/binary"
	"testing"
)

func BenchmarkVarintAppend(b *testing.B) {
	buf := make([]byte, 0, 10)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf = binary.AppendVarint(buf[:0], int64(i))
	}
}

func BenchmarkVarintGet(b *testing.B) {
	buf := make([]byte, 10)
	binary.PutVarint(buf, int64(b.N))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = binary.Varint(buf)
	}
}

func BenchmarkPutInt(b *testing.B) {
	buf := make([]byte, 0, 10)
	var val uint64 = 123

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf = binary.LittleEndian.AppendUint64(buf[:0], val)
	}
}

func BenchmarkPutIntInterface(b *testing.B) {
	buf := make([]byte, 0, 10)
	var val uint64 = 123
	var end binary.AppendByteOrder = binary.LittleEndian

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf = end.AppendUint64(buf[:0], val)
	}
}
