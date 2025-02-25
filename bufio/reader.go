package bufio

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/webmafia/fast"
)

var ErrBufferFull = bufio.ErrBufferFull

var _ BufioReader = (*Reader)(nil)

type Reader struct {
	buf       []byte
	rd        io.Reader // reader provided by the client
	r, w      int       // buf read (left) and write (right) positions
	maxUsed   int       // stats for how much of the buffer was used
	maxSize   int       // buf max size
	totalRead int       // number of total read bytes from io.Reader since last reset
	locked    bool      // locked reader might grow, but never slide
	limited   *LimitedReader
}

// NewReader returns a new [Reader] whose buffer has the default size.
func NewReader(rd io.Reader, bufMaxSize ...int) *Reader {
	if b, ok := rd.(*Reader); ok {
		return b
	}

	r := &Reader{
		rd: rd,
	}

	r.limited = &LimitedReader{
		r: r,
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

func (b *Reader) TotalRead() int {
	return b.totalRead
}

// Size returns the size of the underlying buffer in bytes.
func (b *Reader) Size() int {
	return len(b.buf)
}

func (b *Reader) MaxSize() int {
	return b.maxSize
}

func (b *Reader) MaxUsed() int {
	return b.maxUsed
}

func (b *Reader) SetMaxSize(n int) {
	b.maxSize = n
}

func (b *Reader) ResetSize(size int) {
	if len(b.buf) != size {
		b.buf = fast.MakeNoZero(size)
	}

	b.Reset()
}

func (b *Reader) ResetBytes(data []byte) {
	if br, ok := b.rd.(*bytes.Reader); ok {
		br.Reset(data)
	} else {
		b.rd = bytes.NewReader(data)
	}

	b.Reset()
}

func (b *Reader) ResetReader(r io.Reader) {
	// Avoid infinite recursion.
	if b == r {
		return
	}

	b.rd = r
	b.Reset()
}

func (b *Reader) Reset() {
	b.r = 0
	b.w = 0
	b.maxUsed = 0
	b.totalRead = 0
	b.locked = false
}

// grow copies the buffer to a new, larger buffer so that there are at least n
// bytes of capacity beyond len(b.buf).
func (b *Reader) grow(n int) {
	size := roundPow(b.w + n)

	if b.maxSize > 0 && b.maxSize < size {
		size = b.maxSize
	}

	if size == len(b.buf) {
		return
	}

	// Copy to new buffer. Don't bother with already read data - the old slice
	// won't be GC collected until there are no more references to it.
	buf := fast.MakeNoZero(size)
	b.moveUnreadData(buf)
	b.buf = buf
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
	return false
	return !b.locked && (b.w == b.r || b.r >= 4096 || b.w >= len(b.buf)/2)
}

func (b *Reader) worthGrowing() bool {
	return b.r < 4096 && len(b.buf)-b.w < 4096
}

// fill reads a new chunk into the buffer.
func (b *Reader) fill() (err error) {
	if sizeLeft := len(b.buf) - b.w; sizeLeft < 4096 {
		if b.worthGrowing() || b.locked {
			b.grow(4096)
		} else {
			b.moveUnreadData(b.buf)
		}
	} else if b.worthMoving() {

		// Slide data to the start of buffer
		b.moveUnreadData(b.buf)
	}

	if b.w >= len(b.buf) {
		return ErrBufferFull
	}

	if b.rd == nil {
		panic("b.rd is nil")
	}

	n, err := b.rd.Read(b.buf[b.w:])
	b.totalRead += n
	b.w += n
	b.maxUsed = max(b.maxUsed, b.w)

	return
}

// Peek returns the next n bytes without advancing the reader.
func (b *Reader) Peek(n int) (buf []byte, err error) {
	locked := b.Lock()

	for b.w-b.r < n {
		if err = b.fill(); err != nil {
			break
		}
	}

	if !locked {
		b.Unlock()
	}

	if avail := b.w - b.r; avail < n {
		return
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
			err = b.fill()
		}

		idx := bytes.IndexByte(b.buf[b.r:b.w], c)

		if idx >= 0 {
			discarded += idx
			b.r += idx
			return
		}

		discarded += (b.w - b.r)
		b.r += (b.w - b.r)

		if err != nil {
			break
		}
	}

	return
}

func (b *Reader) ReadSlice(delimiter byte) (slice []byte, err error) {
	locked := b.Lock()
	start := b.r

	for {
		if b.r == b.w {
			err = b.fill()
		}

		idx := bytes.IndexByte(b.buf[b.r:b.w], delimiter)

		if idx >= 0 {
			b.r += idx
			break
		}

		b.r = b.w

		if err != nil {
			break
		}
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

	if err == io.EOF && b.r != b.w {
		err = nil
	}

	return
}

func (b *Reader) ReadByte() (c byte, err error) {
	for b.r == b.w {
		if err = b.fill(); err != nil {
			break
		}
	}

	if avail := b.w - b.r; avail < 1 {
		return
	}

	c = b.buf[b.r]
	b.r++

	return c, nil
}

// Buffered returns the number of bytes that can be read from the current buffer.
func (b *Reader) Buffered() int {
	return b.w - b.r
}

// DebugState returns a string with three lines showing a window of winSize bytes from the
// underlying buffer. The window is chosen so that b.r is as close to the middle of the buffer as possible.
// Line 1 is a marker line: it prints "__" in the cell corresponding to b.r,
// and two spaces for all other cells.
// Line 2 prints each byte as a two-digit hexadecimal number.
// Line 3 prints each byte as a character if printable, or a dot otherwise.
func (b *Reader) DebugState(winSize int) string {
	// Determine the length of the underlying buffer.
	bufLen := len(b.buf)
	if bufLen == 0 || winSize <= 0 {
		return ""
	}
	// Adjust winSize if it exceeds the buffer.
	if winSize > bufLen {
		winSize = bufLen
	}

	// Compute the window start such that b.r is roughly in the middle.
	desiredStart := b.r - winSize/2
	if desiredStart < 0 {
		desiredStart = 0
	}
	if desiredStart+winSize > bufLen {
		desiredStart = bufLen - winSize
	}

	window := b.buf[desiredStart : desiredStart+winSize]
	// Compute pointer position within the window.
	pointerPos := b.r - desiredStart

	// Prepare the three lines.
	markerCells := make([]string, len(window))
	hexCells := make([]string, len(window))
	textCells := make([]string, len(window))
	for i, v := range window {
		// Line 1: Marker – print "__" at pointer position, "  " elsewhere.
		if i == pointerPos {
			markerCells[i] = "__"
		} else {
			markerCells[i] = "  "
		}
		// Line 2: Two-digit hexadecimal representation.
		hexCells[i] = fmt.Sprintf("%02X", v)

		// Line 3: Printable character or dot.
		var ch string
		if v >= 32 && v < 127 {
			ch = string(v)
		} else {
			ch = "."
		}
		textCells[i] = fmt.Sprintf("%-2s", ch)
	}

	// Join the cells with a space between them.
	markerLine := strings.Join(markerCells, " ")
	hexLine := strings.Join(hexCells, " ")
	textLine := strings.Join(textCells, " ")

	return markerLine + "\n" + hexLine + "\n" + textLine
}
