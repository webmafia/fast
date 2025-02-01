package bufio

import (
	"bytes"
	"errors"
	"io"

	"github.com/webmafia/fast"
)

var (
	ErrInvalidUnreadByte = errors.New("bufio: invalid use of UnreadByte")
	ErrInvalidUnreadRune = errors.New("bufio: invalid use of UnreadRune")
	ErrBufferFull        = errors.New("bufio: buffer full")
	ErrNegativeCount     = errors.New("bufio: negative count")
)

type Reader struct {
	buf     []byte
	rd      io.Reader // reader provided by the client
	r, w    int       // buf read (left) and write (right) positions
	maxSize int       // buf max size
	locked  bool      // locked reader might grow, but never slide
}

// NewReader returns a new [Reader] whose buffer has the default size.
func NewReader(rd io.Reader, bufMaxSize ...int) *Reader {
	if b, ok := rd.(*Reader); ok {
		return b
	}

	r := &Reader{
		rd: rd,
	}

	if len(bufMaxSize) > 0 && bufMaxSize[0] > 0 {
		r.maxSize = bufMaxSize[0]
	}

	return r
}

// Locks the buffer so that it won't be slided. Returns whether it was already locked.
func (b *Reader) Lock() (l bool) {
	l = b.locked
	b.locked = true
	return
}

// Unlock the buffer so that it might slide. Returns whether it was locked.
func (b *Reader) Unlock() (l bool) {
	l = b.locked
	b.locked = false
	return
}

// Size returns the size of the underlying buffer in bytes.
func (b *Reader) Size() int {
	return len(b.buf)
}

func (b *Reader) Reset(r io.Reader) {
	// Avoid infinite recursion.
	if b == r {
		return
	}

	b.reset(r)
}

func (b *Reader) reset(r io.Reader) {
	b.rd = r
	b.r = 0
	b.w = 0
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *Reader) grow(n int) error {
	if len(b.buf) == b.maxSize {
		return ErrBufferFull
	}

	size := roundPow(b.w + n)

	if b.maxSize > 0 && b.maxSize < size {
		size = b.maxSize
	}

	// Copy to new buffer. Don't bother with already read data - the old slice
	// won't be GC collected until there are no more references to it.
	buf := fast.MakeNoZero(size)
	b.moveUnreadData(buf)
	b.buf = buf
	return nil
}

func (b *Reader) moveUnreadData(buf []byte) {
	if b.w == b.r {
		b.w = 0
		b.r = 0
	} else {
		copy(buf, b.buf[b.r:b.w])
		b.w -= b.r
		b.r = 0
	}
}

func (b *Reader) worthMoving() bool {
	return !b.locked && (b.w == b.r || b.r >= 4096 || b.w >= len(b.buf)/2)
}

func (b *Reader) worthGrowing() bool {
	return b.r < 4096 && len(b.buf)-b.w < 4096
}

// fill reads a new chunk into the buffer.
func (b *Reader) fill() (err error) {
	if sizeLeft := len(b.buf) - b.w; sizeLeft < 4096 {
		if b.worthGrowing() || b.locked {
			if err = b.grow(4096); err != nil {
				return
			}
		} else {
			b.moveUnreadData(b.buf)
		}
	} else if b.worthMoving() {

		// Slide data to the start of buffer
		b.moveUnreadData(b.buf)
	}

	n, err := b.rd.Read(b.buf[b.w:])
	b.w += n

	return
}

// Peek returns the next n bytes without advancing the reader.
func (b *Reader) Peek(n int) ([]byte, error) {
	locked := b.Lock()

	for b.w-b.r < n {
		if err := b.fill(); err != nil {
			return nil, err
		}
	}

	if !locked {
		b.Unlock()
	}

	return b.buf[b.r : b.r+n], nil
}

func (b *Reader) ReadBytes(n int) (r []byte, err error) {
	r, err = b.Peek(n)
	b.r += len(r)
	return
}

// Discard skips the next n bytes, returning the number of bytes discarded.
func (b *Reader) Discard(n int) (discarded int, err error) {
	if n == 0 {
		return
	}

	remain := n
	for {
		skip := b.Buffered()
		if skip == 0 {
			err = b.fill()
			skip = b.Buffered()
		}
		if skip > remain {
			skip = remain
		}
		b.r += skip
		remain -= skip
		if remain == 0 {
			return n, nil
		}
		if err != nil {
			return n - remain, err
		}
	}
}

func (b *Reader) DiscardUntil(c byte) (discarded int, err error) {
	for {
		if b.r == b.w {
			if err = b.fill(); err != nil {
				return
			}
		}

		idx := bytes.IndexByte(b.buf[b.r:b.w], c)

		if idx >= 0 {
			discarded += idx
			b.r += idx
			return
		}

		discarded += (b.w - b.r)
		b.r += (b.w - b.r)
	}
}

func (b *Reader) ReadSlice(delimiter byte) (slice []byte, err error) {
	locked := b.Lock()
	start := b.r

	for {
		if b.r == b.w {
			if err = b.fill(); err != nil {
				break
			}
		}

		idx := bytes.IndexByte(b.buf[b.r:b.w], delimiter)

		if idx >= 0 {
			b.r += idx
			break
		}

		b.r = b.w
	}

	if !locked {
		b.Unlock()
	}

	slice = b.buf[start:b.r]
	return
}

func (b *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	if buffered := b.w - b.r; buffered < len(p) {
		err = b.fill()
	}

	n = copy(p, b.buf[b.r:b.w])
	b.r += n
	return
}

func (b *Reader) ReadByte() (c byte, err error) {
	for b.r == b.w {
		if err = b.fill(); err != nil {
			return
		}
	}

	c = b.buf[b.r]
	b.r++
	return c, nil
}

// Buffered returns the number of bytes that can be read from the current buffer.
func (b *Reader) Buffered() int {
	return b.w - b.r
}
