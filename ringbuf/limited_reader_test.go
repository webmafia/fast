package ringbuf

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// helperNewLimitedReader creates a LimitedReader wrapping a new Reader with the provided input and limit.
func helperNewLimitedReader(input string, limit int) *LimitedReader {
	r := NewReader(strings.NewReader(input))
	return r.LimitReader(limit)
}

func TestLimitedReaderBuffered(t *testing.T) {
	input := "BufferedTestData"
	limit := len(input) // full length
	lr := helperNewLimitedReader(input, limit)

	// Initially, no data is buffered because nothing has been read.
	if got := lr.Buffered(); got != 0 {
		t.Fatalf("Buffered: got %d, want 0", got)
	}

	// Read a few bytes.
	buf := make([]byte, 5)
	n, err := lr.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	if n != 5 {
		t.Fatalf("Read: expected 5 bytes, got %d", n)
	}

	// Now, some bytes may be buffered in the underlying Reader's ring.
	// Buffered() should report a number up to the remaining limit.
	b := lr.Buffered()
	if b <= 0 || b > limit-5 {
		t.Fatalf("Buffered: got %d, want >0 and <= %d", b, limit-5)
	}
}

func TestLimitedReaderRead(t *testing.T) {
	input := "This is a test for LimitedReader.Read"
	limit := len(input)
	lr := helperNewLimitedReader(input, limit)

	buf := make([]byte, len(input))
	n, err := lr.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	if n != len(input) {
		t.Fatalf("Read: expected %d bytes, got %d", len(input), n)
	}
	if string(buf) != input {
		t.Fatalf("Read: expected %q, got %q", input, string(buf))
	}

	// After reading the full limit, subsequent reads should return EOF.
	_, err = lr.Read(buf)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Read after limit: expected io.EOF, got %v", err)
	}
}

func TestLimitedReaderReadByte(t *testing.T) {
	input := "ABC"
	limit := len(input)
	lr := helperNewLimitedReader(input, limit)

	var out []byte
	for i := 0; i < len(input); i++ {
		b, err := lr.ReadByte()
		if err != nil {
			t.Fatalf("ReadByte error at index %d: %v", i, err)
		}
		out = append(out, b)
	}
	if string(out) != input {
		t.Fatalf("ReadByte: expected %q, got %q", input, string(out))
	}
	// Next ReadByte should return EOF.
	_, err := lr.ReadByte()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("ReadByte after limit: expected io.EOF, got %v", err)
	}
}

func TestLimitedReaderPeek(t *testing.T) {
	input := "PeekLimitedReaderTest"
	limit := len(input)
	lr := helperNewLimitedReader(input, limit)

	// Peek the entire input.
	peek, err := lr.Peek(len(input))
	if err != nil {
		t.Fatalf("Peek error: %v", err)
	}
	if string(peek) != input {
		t.Fatalf("Peek: expected %q, got %q", input, string(peek))
	}
	// Peek should not advance the read pointer.
	if lr.Buffered() != len(input) {
		t.Fatalf("After Peek: expected Buffered() to be %d, got %d", len(input), lr.Buffered())
	}

	// Try peeking more than the limit; should return EOF.
	_, err = lr.Peek(len(input) + 1)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Peek over limit: expected io.EOF, got %v", err)
	}
}

func TestLimitedReaderReadBytes(t *testing.T) {
	input := "ReadBytesLimitedReader"
	limit := len(input)
	lr := helperNewLimitedReader(input, limit)

	b, err := lr.ReadBytes(len(input))
	if err != nil && err != io.EOF {
		t.Fatalf("ReadBytes error: %v", err)
	}
	if string(b) != input {
		t.Fatalf("ReadBytes: expected %q, got %q", input, string(b))
	}
	// Subsequent call should return EOF.
	_, err = lr.ReadBytes(1)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("ReadBytes after limit: expected io.EOF, got %v", err)
	}
}

func TestLimitedReaderDiscard(t *testing.T) {
	input := "DiscardLimitedReaderTest"
	limit := len(input)
	lr := helperNewLimitedReader(input, limit)

	// Discard the first 7 bytes.
	n, err := lr.Discard(7)
	if err != nil && err != io.EOF {
		t.Fatalf("Discard error: %v", err)
	}
	if n != 7 {
		t.Fatalf("Discard: expected to discard 7 bytes, got %d", n)
	}
	// Read the rest.
	rest, err := io.ReadAll(lr)
	if err != nil && err != io.EOF {
		t.Fatalf("ReadAll error after Discard: %v", err)
	}
	expected := input[7:]
	if string(rest) != expected {
		t.Fatalf("After Discard: expected %q, got %q", expected, string(rest))
	}
}

func TestLimitedReaderDiscardUntil(t *testing.T) {
	input := "abcXYZdef"
	limit := len(input)
	lr := helperNewLimitedReader(input, limit)

	// Discard until we hit 'X'.
	n, err := lr.DiscardUntil('X')
	if err != nil && err != io.EOF {
		t.Fatalf("DiscardUntil error: %v", err)
	}
	// Expect to discard "abc" (3 bytes).
	if n != 3 {
		t.Fatalf("DiscardUntil: expected to discard 3 bytes, got %d", n)
	}
	// Next byte should be 'X'.
	b, err := lr.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte after DiscardUntil error: %v", err)
	}
	if b != 'X' {
		t.Fatalf("DiscardUntil: expected next byte 'X', got %q", b)
	}
}

func TestLimitedReaderPartialFill(t *testing.T) {
	// Use a slowReader that provides data in fixed-size chunks.
	input := "PartialFillLimitedReaderTestData"
	limit := len(input)
	custom := &slowReader{data: []byte(input), chunk: 5}
	// Wrap the slowReader with our Reader and then a LimitedReader.
	r := NewReader(custom)
	lr := r.LimitReader(limit)

	buf := make([]byte, len(input))
	n, err := lr.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Partial fill Read error: %v", err)
	}
	if n != len(input) {
		t.Fatalf("Partial fill Read: expected %d bytes, got %d", len(input), n)
	}
	if string(buf) != input {
		t.Fatalf("Partial fill Read: expected %q, got %q", input, string(buf))
	}
}

func TestLimitedReaderLimitExhaustion(t *testing.T) {
	input := "LimitExhaustionTest"
	limit := 10 // Set limit to fewer bytes than input length.
	lr := helperNewLimitedReader(input, limit)

	// Read until limit is reached.
	buf := make([]byte, 100)
	n, err := lr.Read(buf)
	if err != nil && err != io.EOF {
		t.Fatalf("Read error: %v", err)
	}
	if n != limit {
		t.Fatalf("Limit exhaustion: expected %d bytes read, got %d", limit, n)
	}
	// Subsequent reads must return EOF.
	n, err = lr.Read(buf)
	if !errors.Is(err, io.EOF) {
		t.Fatalf("After limit reached: expected io.EOF, got %v", err)
	}
	if n != 0 {
		t.Fatalf("After limit reached: expected 0 bytes, got %d", n)
	}
}
