package ringbuf

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// TestReaderBuffered verifies that Buffered returns the number of unread bytes.
func TestReaderBuffered(t *testing.T) {
	input := "Hello, ringbuf!"
	src := strings.NewReader(input)
	r := NewReader(src)
	// Initially, nothing is buffered.
	if got := r.Buffered(); got != 0 {
		t.Fatalf("Buffered: got %d, want 0", got)
	}
	// Trigger a fill by reading a few bytes.
	buf := make([]byte, 5)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	// After reading 5 bytes, some bytes may still be buffered.
	if got := r.Buffered(); got == 0 {
		t.Fatalf("Buffered: got %d, want >0", got)
	}
	if n != 5 {
		t.Fatalf("Read: expected 5 bytes, got %d", n)
	}
}

// TestReaderRead reads all bytes from the Reader and compares with the input.
func TestReaderRead(t *testing.T) {
	input := "This is a test of the Reader."
	src := strings.NewReader(input)
	r := NewReader(src)
	buf := make([]byte, len(input))
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	if n != len(input) {
		t.Fatalf("Read: expected %d bytes, got %d", len(input), n)
	}
	if string(buf) != input {
		t.Fatalf("Read: expected %q, got %q", input, string(buf))
	}
}

// TestReaderReadByte reads one byte at a time.
func TestReaderReadByte(t *testing.T) {
	input := "ABC"
	src := strings.NewReader(input)
	r := NewReader(src)
	var out []byte
	for i := 0; i < len(input); i++ {
		b, err := r.ReadByte()
		if err != nil {
			t.Fatalf("ReadByte error at %d: %v", i, err)
		}
		out = append(out, b)
	}
	if string(out) != input {
		t.Fatalf("ReadByte: expected %q, got %q", input, string(out))
	}
	// Next ReadByte should return EOF.
	_, err := r.ReadByte()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("ReadByte after data: expected EOF, got %v", err)
	}
}

// TestReaderPeek verifies that Peek returns the next n bytes without advancing.
func TestReaderPeek(t *testing.T) {
	input := "PeekTestData"
	src := strings.NewReader(input)
	r := NewReader(src)
	peek, err := r.Peek(len(input))
	if err != nil {
		t.Fatalf("Peek error: %v", err)
	}
	if string(peek) != input {
		t.Fatalf("Peek: expected %q, got %q", input, string(peek))
	}
	// Now, read the data and verify it matches.
	buf := make([]byte, len(input))
	_, err = r.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	if string(buf) != input {
		t.Fatalf("Read after Peek: expected %q, got %q", input, string(buf))
	}
}

// TestReaderReadBytes reads exactly n bytes using ReadBytes.
func TestReaderReadBytes(t *testing.T) {
	input := "ReadBytesTest"
	src := strings.NewReader(input)
	r := NewReader(src)
	b, err := r.ReadBytes(len(input))
	if err != nil && err != io.EOF {
		t.Fatalf("ReadBytes error: %v", err)
	}
	if string(b) != input {
		t.Fatalf("ReadBytes: expected %q, got %q", input, string(b))
	}
	// Subsequent ReadBytes should return EOF.
	_, err = r.ReadBytes(1)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("ReadBytes after EOF: expected EOF, got %v", err)
	}
}

// TestReaderDiscard discards a specified number of bytes.
func TestReaderDiscard(t *testing.T) {
	input := "DiscardTestData"
	src := strings.NewReader(input)
	r := NewReader(src)
	// Discard the first 7 bytes ("Discard")
	n, err := r.Discard(7)
	if err != nil && err != io.EOF {
		t.Fatalf("Discard error: %v", err)
	}
	if n != 7 {
		t.Fatalf("Discard: expected 7, got %d", n)
	}
	// Read the remaining data.
	rest, err := io.ReadAll(r)
	if err != nil && err != io.EOF {
		t.Fatalf("ReadAll error: %v", err)
	}
	expected := "TestData"
	if string(rest) != expected {
		t.Fatalf("After Discard: expected %q, got %q", expected, string(rest))
	}
}

// TestReaderDiscardUntil discards bytes until a specific byte is encountered.
func TestReaderDiscardUntil(t *testing.T) {
	input := "abcXYZdef"
	src := strings.NewReader(input)
	r := NewReader(src)
	// Discard until 'X' is found.
	n, err := r.DiscardUntil('X')
	if err != nil && err != io.EOF {
		t.Fatalf("DiscardUntil error: %v", err)
	}
	// Expect to discard "abc" (3 bytes).
	if n != 3 {
		t.Fatalf("DiscardUntil: expected to discard 3 bytes, got %d", n)
	}
	// Next byte should be 'X'.
	b, err := r.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte error after DiscardUntil: %v", err)
	}
	if b != 'X' {
		t.Fatalf("DiscardUntil: expected next byte 'X', got %q", b)
	}
}

// TestReaderPartialFill simulates an underlying reader that provides data in small chunks.
func TestReaderPartialFill(t *testing.T) {
	input := "PartialFillTestData"
	custom := &slowReader{data: []byte(input), chunk: 5}
	r := NewReader(custom)
	buf := make([]byte, len(input))
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	if n != len(input) {
		t.Fatalf("Read: expected %d bytes, got %d", len(input), n)
	}
	if string(buf) != input {
		t.Fatalf("Read: expected %q, got %q", input, string(buf))
	}
}

// slowReader simulates a reader that returns data in fixed-size chunks.
type slowReader struct {
	data  []byte
	chunk int
	pos   int
}

func (s *slowReader) Read(p []byte) (int, error) {
	if s.pos >= len(s.data) {
		return 0, io.EOF
	}
	n := s.chunk
	if s.pos+n > len(s.data) {
		n = len(s.data) - s.pos
	}
	copy(p, s.data[s.pos:s.pos+n])
	s.pos += n
	return n, nil
}

// TestReaderPeekWrap verifies Peek when data wraps around in the RingBuf.
// We simulate wrap-around by reading enough bytes from the underlying reader.
func TestReaderPeekWrap(t *testing.T) {
	// Create input longer than BufferSize so that wrap-around occurs.
	// We'll write BufferSize-10 "A"s, then force a wrap by reading some bytes and then writing "B"s.
	pattern := strings.Repeat("A", BufferSize-10)
	extra := strings.Repeat("B", 20)
	input := pattern + extra
	src := strings.NewReader(input)
	r := NewReader(src)
	// Read enough to force the ring's read pointer near the end.
	// We'll read (BufferSize - 10) bytes from the underlying reader.
	buf := make([]byte, BufferSize-10)
	if _, err := io.ReadFull(r, buf); err != nil {
		t.Fatalf("Pre-read error: %v", err)
	}
	// Now, the ring buffer should have wrapped when filling extra.
	// Peek the extra 20 bytes.
	peek, err := r.Peek(20)
	if err != nil {
		t.Fatalf("Peek error: %v", err)
	}
	if string(peek) != extra {
		t.Fatalf("Peek wrap: expected %q, got %q", extra, string(peek))
	}
}
