package buffer

import (
	"testing"
)

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
