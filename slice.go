package fast

import "unsafe"

//go:inline
func SliceToSlice[From any, To any](from []From, toLength int) []To {
	fromHeader := (*sliceHeader)(unsafe.Pointer(&from))
	toHeader := sliceHeader{
		Data: fromHeader.Data,
		Len:  toLength,
		Cap:  toLength,
	}

	return *(*[]To)(unsafe.Pointer(&toHeader))
}

//go:inline
func SliceUnsafePointer[T any](slice []T) unsafe.Pointer {
	header := *(*sliceHeader)(unsafe.Pointer(&slice))
	return header.Data
}
