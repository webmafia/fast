package buffer

import (
	"fmt"
	"testing"
)

// From: https://github.com/valyala/bytebufferpool
func originalIndex(n int) int {
	n--
	n >>= minBitSize
	idx := 0
	for n > 0 {
		n >>= 1
		idx++
	}
	if idx >= steps {
		idx = steps - 1
	}
	return idx
}

func TestIndex(t *testing.T) {
	for i := range maxSize {
		a := index(i)
		b := originalIndex(i)

		if a != b {
			t.Fatalf("%d: expected %d, got %d", i, b, a)
		}
	}
}

func Example_index() {
	for _, i := range [...]int{64, 128, 192, 256, 257} {
		fmt.Println(i, "=", index(i))
	}

	// Output:
	//
	// 64 = 0
	// 128 = 1
	// 192 = 2
	// 256 = 2
	// 257 = 3
}

func BenchmarkIndex(b *testing.B) {
	b.Run("Original", func(b *testing.B) {
		for i := range b.N {
			_ = originalIndex(i)
		}
	})

	b.Run("Modified", func(b *testing.B) {
		for i := range b.N {
			_ = index(i)
		}
	})
}

func Example_roundPow() {
	for _, i := range [...]int{0, 1, 2, 64, 128, 192, 256, 257} {
		fmt.Println(i, "=", roundPow(i))
	}

	// Output:
	//
	// 0 = 1
	// 1 = 1
	// 2 = 2
	// 64 = 64
	// 128 = 128
	// 192 = 256
	// 256 = 256
	// 257 = 512
}

func Test_roundPow(t *testing.T) {
	for i := 64; i <= maxSize; i++ {
		rounded := roundPow(i)
		idx := index(i)
		expected := (1 << (idx + minBitSize))

		if rounded != expected {
			t.Fatalf("%d: expected %d, got %d", i, expected, rounded)
		}
	}
}

func Benchmark_roundPow(b *testing.B) {
	for i := range b.N {
		_ = roundPow(i)
	}
}
