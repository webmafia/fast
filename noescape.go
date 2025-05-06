package fast

import "unsafe"

// noescape hides a pointer from escape analysis. It is the identity function
// but escape analysis doesn't think the output depends on the input.
// noescape is inlined and currently compiles down to zero instructions.
// USE CAREFULLY!
// This was copied from the runtime; see issues 23382 and 7921.
//
//go:nosplit
//go:nocheckptr
func NoescapeUnsafe(p unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) ^ 0)
}

//go:inline
func Noescape[T any](p T) T {
	return *(*T)(NoescapeUnsafe(unsafe.Pointer(&p)))
}
