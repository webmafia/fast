package fast

import (
	"fmt"
	"testing"
)

func ExampleMakeNoZero() {
	b := MakeNoZero(16)
	fmt.Printf("Length %d, capacity %d", len(b), cap(b))
	// Output: Length 16, capacity 16
}

func stdMake(n int) []byte {
	return make([]byte, n)
}

func BenchmarkMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = stdMake(256)
	}
}

func BenchmarkMakeNoZero(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MakeNoZero(256)
	}
}
