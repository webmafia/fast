package fast

// See: https://pkg.go.dev/encoding@master#pkg-types

type TextAppender interface {
	// AppendText appends the textual representation of itself to the end of b
	// (allocating a larger slice if necessary) and returns the updated slice.
	//
	// Implementations must not retain b, nor mutate any bytes within b[:len(b)].
	AppendText(b []byte) ([]byte, error)
}

type BinaryAppender interface {
	// AppendBinary appends the binary representation of itself to the end of b
	// (allocating a larger slice if necessary) and returns the updated slice.
	//
	// Implementations must not retain b, nor mutate any bytes within b[:len(b)].
	AppendBinary(b []byte) ([]byte, error)
}

type JsonAppender interface {
	// AppendJson appends the JSON representation of itself to the end of b
	// (allocating a larger slice if necessary) and returns the updated slice.
	//
	// Implementations must not retain b, nor mutate any bytes within b[:len(b)].
	AppendJson(b []byte) ([]byte, error)
}
