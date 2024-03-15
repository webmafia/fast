package fast

import (
	"sync"
	"testing"
)

func BenchmarkPool(b *testing.B) {
	pool := sync.Pool{
		New: func() any {
			return &struct{}{}
		},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := pool.Get()
		pool.Put(v)
	}
}

func BenchmarkGenericPool(b *testing.B) {
	pool := NewPool[struct{}]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := pool.Acquire()
		pool.Release(v)
	}
}

func BenchmarkPool_Parallell(b *testing.B) {
	pool := sync.Pool{
		New: func() any {
			return &struct{}{}
		},
	}

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v := pool.Get()
			pool.Put(v)
		}
	})
}

func BenchmarkGenericPool_Parallell(b *testing.B) {
	pool := NewPool[struct{}]()
	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			v := pool.Acquire()
			pool.Release(v)
		}
	})
}
