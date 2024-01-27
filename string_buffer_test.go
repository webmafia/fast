package fast

import (
	"fmt"
	"testing"
)

func BenchmarkStringBuffer_WriteFloat64(b *testing.B) {
	var buf StringBuffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.WriteFloat64(123.456789)
		buf.Reset()
	}
}

func BenchmarkStringBuffer_WriteFloat64Lossy(b *testing.B) {
	var buf StringBuffer

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf.WriteFloat64Lossy(123.456789)
		buf.Reset()
	}
}

func ExampleStringBuffer() {
	var b StringBuffer
	b.WriteString("hello")
	b.WriteByte(' ')
	b.WriteInt(123)
	fmt.Println(b)

	b.Reset()

	b.WriteFloat64Lossy(123.4567891)
	b.WriteByte(' ')
	b.WriteBool(true)
	fmt.Println(b)

	// Output:
	// hello 123
	// 123.456789 true
}
