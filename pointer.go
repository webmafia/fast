package fast

import "unsafe"

//go:inline
func PointerToBytes[T any](val *T, length int) []byte {
	header := sliceHeader{
		Data: unsafe.Pointer(val),
		Len:  length,
		Cap:  length,
	}

	return *(*[]byte)(unsafe.Pointer(&header))
}

//go:inline
func PointerToBytesOffset[T any](val *T, length, offset int) []byte {
	header := sliceHeader{
		Data: unsafe.Add(unsafe.Pointer(val), offset),
		Len:  length,
		Cap:  length,
	}

	return *(*[]byte)(unsafe.Pointer(&header))
}

//go:inline
func BytesToPointer[T any](b []byte) *T {
	header := *(*sliceHeader)(unsafe.Pointer(&b))
	return (*T)(header.Data)
}
