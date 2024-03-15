package fast

import (
	"sync/atomic"
	"unsafe"
)

type AtomicInt32Pair struct {
	v int64
}

type atomicInt32Pair struct {
	a, b int32
}

//go:inline
func toInt32(v int64) (int32, int32) {
	ab := *(*atomicInt32Pair)(unsafe.Pointer(&v))
	return ab.a, ab.b
}

//go:inline
func toInt64(a, b int32) int64 {
	return *(*int64)(unsafe.Pointer(&atomicInt32Pair{
		a: a,
		b: b,
	}))
}

// Load atomically loads and returns the value stored in x.
func (x *AtomicInt32Pair) Load() (int32, int32) {
	return toInt32(atomic.LoadInt64(&x.v))
}

// Store atomically stores val into x.
func (x *AtomicInt32Pair) Store(a, b int32) {
	atomic.StoreInt64(&x.v, toInt64(a, b))
}

// Swap atomically stores new into x and returns the previous value.
func (x *AtomicInt32Pair) Swap(newA, newB int32) (oldA, oldB int32) {
	return toInt32(atomic.SwapInt64(&x.v, toInt64(newA, newB)))
}

// CompareAndSwap executes the compare-and-swap operation for x.
func (x *AtomicInt32Pair) CompareAndSwap(oldA, oldB, newA, newB int32) (swapped bool) {
	return atomic.CompareAndSwapInt64(&x.v, toInt64(oldA, oldB), toInt64(newA, newB))
}

// Add atomically adds delta to x and returns the new value.
func (x *AtomicInt32Pair) Add(deltaA, deltaB int32) (newA, newB int32) {
	return toInt32(atomic.AddInt64(&x.v, toInt64(deltaA, deltaB)))
}
