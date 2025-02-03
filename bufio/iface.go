package bufio

import "io"

type BufioReader interface {
	io.Reader
	io.ByteReader
	Discard(n int) (discarded int, err error)
	DiscardUntil(c byte) (discarded int, err error)
	Peek(n int) (buf []byte, err error)
	ReadBytes(n int) (r []byte, err error)
	ReadSlice(delimiter byte) (slice []byte, err error)
	Buffered() int
	Lock() bool
	Unlock() bool
	Size() int
}
