package fast

import (
	_ "runtime"
	"unsafe"
)

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

// MakeNoZero makes a slice of length and capacity n without zeroing the bytes.
// It is the caller's responsibility to ensure uninitialized bytes
// do not leak to the end user.
func MakeNoZero(len int) []byte {
	return unsafe.Slice((*byte)(mallocgc(uintptr(len), nil, false)), len)
}
