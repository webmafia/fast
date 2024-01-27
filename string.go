package fast

import "unsafe"

//go:inline
func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

//go:inline
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
