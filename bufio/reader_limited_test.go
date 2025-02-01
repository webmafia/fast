package bufio

import (
	"bytes"
	"io"
	"testing"
)

func TestLimitedReader_Read(t *testing.T) {
	data := []byte("abcdef")
	reader := NewReader(bytes.NewReader(data))
	limited := reader.LimitReader(3)

	buf := make([]byte, 4)
	n, err := limited.Read(buf)
	if n != 3 || err != io.EOF && err != nil {
		t.Errorf("LimitedReader.Read() = %d, %v; want 3, io.EOF", n, err)
	}
}

func TestLimitedReader_ReadByte(t *testing.T) {
	data := []byte("hello")
	reader := NewReader(bytes.NewReader(data))
	limited := reader.LimitReader(2)

	b, err := limited.ReadByte()
	if b != 'h' || err != nil {
		t.Errorf("LimitedReader.ReadByte() = %q, %v; want 'h', nil", b, err)
	}

	b, err = limited.ReadByte()
	if b != 'e' || err != nil {
		t.Errorf("LimitedReader.ReadByte() = %q, %v; want 'e', nil", b, err)
	}

	_, err = limited.ReadByte()
	if err != io.EOF {
		t.Errorf("LimitedReader.ReadByte() beyond limit should return io.EOF, got %v", err)
	}
}

func TestLimitedReader_Discard(t *testing.T) {
	data := []byte("abcdef")
	reader := NewReader(bytes.NewReader(data))
	limited := reader.LimitReader(5)

	discarded, err := limited.Discard(3)
	if discarded != 3 || err != nil {
		t.Errorf("LimitedReader.Discard() = %d, %v; want 3, nil", discarded, err)
	}

	discarded, err = limited.Discard(3)
	if discarded != 2 || err != io.EOF {
		t.Errorf("LimitedReader.Discard() = %d, %v; want 2, io.EOF", discarded, err)
	}
}

func TestLimitedReader_Peek(t *testing.T) {
	data := []byte("abcdef")
	reader := NewReader(bytes.NewReader(data))
	limited := reader.LimitReader(4)

	buf, err := limited.Peek(2)
	if err != nil {
		t.Errorf("LimitedReader.Peek() failed: %v", err)
	}
	if !bytes.Equal(buf, []byte("ab")) {
		t.Errorf("LimitedReader.Peek() = %q; want 'ab'", buf)
	}

	_, err = limited.Peek(5)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("LimitedReader.Peek() beyond limit should return io.ErrUnexpectedEOF, got %v", err)
	}
}

func TestLimitedReader_Buffered(t *testing.T) {
	data := []byte("hello")
	reader := NewReader(bytes.NewReader(data))
	limited := reader.LimitReader(3)
	reader.fill()

	if buffered := limited.Buffered(); buffered > 3 {
		t.Errorf("LimitedReader.Buffered() = %d; should not exceed limit of 3", buffered)
	}
}
