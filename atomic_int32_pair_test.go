package fast

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func ExampleAtomicInt32Pair() {
	var v AtomicInt32Pair

	for i := 0; i < 100; i++ {
		v.Add(1, -1)
	}

	fmt.Println(v.Load())

	// Output: 100 -100
}

func BenchmarkAtomicInt32Pair_Add(b *testing.B) {
	var v AtomicInt32Pair

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v.Add(1, -1)
	}
}

func BenchmarkAtomicInt64_Add(b *testing.B) {
	var v atomic.Int64

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v.Add(1)
	}
}

func BenchmarkAtomicInt32_Add(b *testing.B) {
	var v atomic.Int32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v.Add(1)
	}
}

func BenchmarkAtomicInt32Pair_Add_Parallell(b *testing.B) {
	var v AtomicInt32Pair

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v.Add(1, -1)
		}
	})
}

func BenchmarkAtomicInt64_Add_Parallell(b *testing.B) {
	var v atomic.Int64

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v.Add(1)
		}
	})
}

func BenchmarkAtomicInt32_Add_Parallell(b *testing.B) {
	var v atomic.Int32

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v.Add(1)
		}
	})
}
