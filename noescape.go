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
func Noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

//go:inline
func NoescapeVal[T any](p *T) *T {
	return (*T)(Noescape(unsafe.Pointer(p)))
}

//go:inline
func NoescapeBytes(b []byte) []byte {
	return *(*[]byte)(Noescape(unsafe.Pointer(&b)))
}
