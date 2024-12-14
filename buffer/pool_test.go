package buffer

import (
	"fmt"
	"testing"
)

func Example_index() {
	for _, i := range [...]int{64} {
		fmt.Println(i, "=", index(i))
	}

	// Output: TODO
}

func Benchmark_calibrate(b *testing.B) {
	var p Pool

	b.ResetTimer()

	for range b.N {
		p.calibrate()
	}
}

func BenchmarkPool(b *testing.B) {
	var p Pool
	data := make([]byte, 32)

	for range b.N {
		buf := p.Get()
		buf.B = append(buf.B, data...)
		p.Put(buf)
	}
}

func BenchmarkPool_Parallell(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		var pool Pool
		data := make([]byte, 32)

		for p.Next() {
			buf := pool.Get()
			buf.B = append(buf.B, data...)
			pool.Put(buf)
		}
	})
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
