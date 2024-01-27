package fast

import "sync"

type Pool[T any] struct {
	pool  sync.Pool
	init  func(*T)
	reset func(*T)
}

// Created a new pool and accepts two (2) optional callbacks. The first is a initializer, and will be called
// whenever a new item (T) is created. The last is a resetter, and will be called whenever an item is
// released back to the pool.
func NewPool[T any](cb ...func(*T)) *Pool[T] {
	p := new(Pool[T])

	if len(cb) > 0 {
		p.init = cb[0]
	}

	if len(cb) > 1 {
		p.reset = cb[1]
	}

	return p
}

// Acquires an item from the pool.
func (p *Pool[T]) Acquire() *T {
	if v, ok := p.pool.Get().(*T); ok {
		return v
	}

	v := new(T)

	if p.init != nil {
		p.init(v)
	}

	return v
}

// Releases an item back to the pool. The item cannot be used after release.
func (p *Pool[T]) Release(v *T) {
	if p.reset != nil {
		p.reset(v)
	}

	p.pool.Put(v)
}
