package bufio

import "io"

type chunkedReader struct {
	data  []byte
	pos   int
	chunk int
}

// NewChunkedReader initializes a new ChunkedReader with a fixed chunk size.
func newChunkedReader(data []byte, chunkSize int) *chunkedReader {
	if chunkSize <= 0 {
		panic("invalid chunk size")
	}
	return &chunkedReader{
		data:  data,
		chunk: chunkSize,
	}
}

func (c *chunkedReader) Read(p []byte) (n int, err error) {
	p = p[:min(len(p), c.chunk, len(c.data)-c.pos)]
	n = copy(p, c.data[c.pos:])
	c.pos += n

	if c.pos >= len(c.data) {
		err = io.EOF
	}

	return
}

func (c *chunkedReader) Reset(data []byte) {
	c.data = data
	c.pos = 0
}
