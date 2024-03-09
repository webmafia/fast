package fast

import (
	_ "runtime"
	"unsafe"
)

// MakeNoZero makes a slice of length and capacity l without zeroing the bytes.
// It is the caller's responsibility to ensure uninitialized bytes
// do not leak to the end user.
func MakeNoZero(l int) []byte {
	return unsafe.Slice((*byte)(mallocgc(uintptr(l), nil, false)), l)
}

// MakeNoZero makes a slice of length l and capacity c without zeroing the bytes.
// It is the caller's responsibility to ensure uninitialized bytes
// do not leak to the end user.
func MakeNoZeroCap(l int, c int) []byte {
	return MakeNoZero(c)[:l]
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer
