package fast

import (
	"sync"
	"testing"
)

func BenchmarkStringBufferPool(b *testing.B) {
	pool := NewStringBufferPool()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := pool.Acquire()
		pool.Release(v)
	}
}

func BenchmarkStringBufferAnyPool(b *testing.B) {
	pool := sync.Pool{
		New: func() any {
			return new(StringBuffer)
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := pool.Get().(*StringBuffer)
		v.Reset()
		pool.Put(v)
	}
}
