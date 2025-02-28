package ringbuf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func Example() {
	var rb RingBuf

	buf := make([]byte, 1000)

	fmt.Println("writing 5 x 1000 bytes")
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))
	fmt.Println(rb.Write(buf))

	fmt.Println("reading 1000 and then writing 1000 more")
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Write(buf))

	fmt.Println("reading 5 x 1000 bytes")
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))
	fmt.Println(rb.Read(buf))

	fmt.Println("reading one additional time")
	fmt.Println(rb.Read(buf))

	// Output:
	//
	// writing 5 x 1000 bytes
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 96 <nil>
	// reading 1000 and then writing 1000 more
	// 1000 <nil>
	// 1000 <nil>
	// reading 5 x 1000 bytes
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 1000 <nil>
	// 96 <nil>
	// reading one additional time
	// 0 EOF
}

// func ExampleRingBuf_DebugDump() {
// 	var r RingBuf

// 	r.Write([]byte("hello"))

// 	r.DebugDump(os.Stdout)

// 	// Output: TODO
// }

func TestWriteAndRead(t *testing.T) {
	var rb RingBuf
	data := []byte("hello, ringbuf!")
	// Write data.
	n, err := rb.Write(data)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Write: got %d bytes, want %d", n, len(data))
	}

	// Read data back.
	buf := make([]byte, len(data))
	n, err = rb.Read(buf)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Read: got %d bytes, want %d", n, len(data))
	}
	if !bytes.Equal(buf, data) {
		t.Fatalf("Read: got %q, want %q", buf, data)
	}

	// Further read should return EOF.
	n, err = rb.Read(buf)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Read after empty: expected io.EOF, got err=%v", err)
	}
	if n != 0 {
		t.Fatalf("Read after empty: got %d bytes, want 0", n)
	}
}

func TestReadByte(t *testing.T) {
	var rb RingBuf
	data := []byte("ABC")
	if n, err := rb.Write(data); err != nil || n != len(data) {
		t.Fatalf("Write error: n=%d, err=%v", n, err)
	}

	var out []byte
	for i := 0; i < len(data); i++ {
		b, err := rb.ReadByte()
		if err != nil {
			t.Fatalf("ReadByte error at %d: %v", i, err)
		}
		out = append(out, b)
	}
	if !bytes.Equal(out, data) {
		t.Fatalf("ReadByte: got %q, want %q", out, data)
	}

	_, err := rb.ReadByte()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("ReadByte on empty: expected io.EOF, got %v", err)
	}
}

func TestPeekContiguous(t *testing.T) {
	var rb RingBuf
	data := []byte(strings.Repeat("X", 100))
	if n, err := rb.Write(data); err != nil || n != len(data) {
		t.Fatalf("Write error: n=%d, err=%v", n, err)
	}
	peekLen := 50
	peek, err := rb.Peek(peekLen)
	if err != nil {
		t.Fatalf("Peek error: %v", err)
	}
	if len(peek) != peekLen {
		t.Fatalf("Peek: got length %d, want %d", len(peek), peekLen)
	}
	if !bytes.Equal(peek, data[:peekLen]) {
		t.Fatalf("Peek: got %q, want %q", peek, data[:peekLen])
	}
	// Ensure Peek does not advance the read pointer.
	if rb.read != 0 {
		t.Fatalf("Peek advanced read pointer: got %d, want 0", rb.read)
	}
}

func TestPeekWrap(t *testing.T) {
	var rb RingBuf
	// Fill the buffer with a pattern.
	pattern := []byte(strings.Repeat("A", int(BufferSize)))
	if n, err := rb.Write(pattern); err != nil || n != int(BufferSize) {
		t.Fatalf("Write error: n=%d, err=%v", n, err)
	}
	// Manually adjust read pointer so that it is near the end.
	_, _ = rb.ReadBytes(BufferSize - 10)

	// At this point, unread() should be: (BufferSize - (BufferSize-10)) = 10.
	// Now write extra data so unread becomes 10 + 20 = 30.
	extra := []byte(strings.Repeat("B", 20))
	if n, err := rb.Write(extra); err != nil || n != len(extra) {
		t.Fatalf("Write extra error: n=%d, err=%v", n, err)
	}
	if rb.unread() != 30 {
		t.Fatalf("Unread: got %d, want %d", rb.unread(), 30)
	}
	peek, err := rb.Peek(30)
	if err != nil {
		t.Fatalf("Peek error: %v", err)
	}
	if len(peek) != 30 {
		t.Fatalf("Peek: got length %d, want %d", len(peek), 30)
	}
	expected := append(pattern[BufferSize-10:], extra...)
	if !bytes.Equal(peek, expected) {
		t.Fatalf("Peek wrap: got %q, want %q", peek, expected)
	}
}

func TestFillFrom(t *testing.T) {
	var rb RingBuf
	srcStr := "HelloFillFrom"
	src := strings.NewReader(srcStr)
	n, err := rb.FillFrom(src)
	if err != nil && err != io.EOF {
		t.Fatalf("FillFrom error: %v", err)
	}
	if n != int64(len(srcStr)) {
		t.Fatalf("FillFrom: got %d, want %d", n, len(srcStr))
	}

	buf := make([]byte, n)
	if rn, err := rb.Read(buf); err != nil {
		t.Fatalf("Read after FillFrom error: %v", err)
	} else if rn != int(n) {
		t.Fatalf("Read after FillFrom: got %d, want %d", rn, n)
	}
	if !bytes.Equal(buf, []byte(srcStr)) {
		t.Fatalf("FillFrom data mismatch: got %q", buf)
	}
}

func TestReadFrom(t *testing.T) {
	var rb RingBuf
	src := strings.NewReader(strings.Repeat("X", int(BufferSize)))
	n, err := rb.ReadFrom(src)
	if err != nil {
		t.Fatalf("ReadFrom error: %v", err)
	}
	if n != int64(BufferSize) {
		t.Fatalf("ReadFrom: got %d, want %d", n, BufferSize)
	}
	buf := make([]byte, BufferSize)
	rn, err := rb.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("Read error: %v", err)
	}
	if rn != BufferSize {
		t.Fatalf("Read: got %d, want %d", rn, BufferSize)
	}
	for i, b := range buf {
		if b != 'X' {
			t.Fatalf("Byte %d: got %q, want 'X'", i, b)
		}
	}
}

func TestWriteError(t *testing.T) {
	var rb RingBuf
	data := make([]byte, BufferSize)
	for i := range data {
		data[i] = 'Z'
	}
	n, err := rb.Write(data)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if n != BufferSize {
		t.Fatalf("Write: wrote %d bytes, want %d", n, BufferSize)
	}
	n, err = rb.Write([]byte("A"))
	if err != io.ErrShortBuffer {
		t.Fatalf("Write extra: expected io.ErrShortBuffer, got %v", err)
	}
	if n != 0 {
		t.Fatalf("Write extra: wrote %d bytes, want 0", n)
	}
}

func TestReadEmpty(t *testing.T) {
	var rb RingBuf
	buf := make([]byte, 10)
	n, err := rb.Read(buf)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Read empty: expected io.EOF, got err=%v", err)
	}
	if n != 0 {
		t.Fatalf("Read empty: got %d bytes, want 0", n)
	}
}
