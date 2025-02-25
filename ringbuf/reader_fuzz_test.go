package ringbuf

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// FuzzReaderOps performs a fuzz test on the Reader by simulating a sequence of
// randomized operations and comparing the output with the expected data.
func FuzzReaderOps(f *testing.F) {
	// Add some seed corpus.
	f.Add("The quick brown fox jumps over the lazy dog.")
	f.Add("Lorem ipsum dolor sit amet, consectetur adipiscing elit.")

	f.Fuzz(func(t *testing.T, data string) {
		// Use the input string as our golden data.
		expected := []byte(data)
		// Create a new ring-buffered Reader.
		r := NewReader(strings.NewReader(data))
		// offset tracks how many bytes we expect have been consumed.
		offset := 0

		// Run a fixed number of randomized operations (or until we've consumed all data).
		ops := 1000
		for i := 0; i < ops && offset < len(expected); i++ {
			// Choose an operation using one byte from data.
			op := int(expected[i%len(expected)]) % 6

			switch op {
			case 0: // Read
				// Pick a random read size between 1 and (remaining+1).
				maxRead := len(expected) - offset
				if maxRead == 0 {
					break
				}
				n := int(expected[(i+1)%len(expected)])%(maxRead) + 1
				buf := make([]byte, n)
				nr, err := r.Read(buf)
				// If we get an error other than EOF before consuming some data, report it.
				if err != nil && nr == 0 && offset < len(expected) {
					t.Fatalf("Read error: %v", err)
				}
				// Compare the data read with the expected slice.
				if !bytes.Equal(buf[:nr], expected[offset:offset+nr]) {
					t.Errorf("Read: got %q, want %q", buf[:nr], expected[offset:offset+nr])
				}
				offset += nr

			case 1: // ReadByte
				if offset >= len(expected) {
					break
				}
				b, err := r.ReadByte()
				if err != nil {
					if offset < len(expected) {
						t.Fatalf("ReadByte error: %v", err)
					}
				} else {
					if b != expected[offset] {
						t.Errorf("ReadByte: got %q, want %q", b, expected[offset])
					}
					offset++
				}

			case 2: // Peek
				// Randomly choose a peek length between 1 and (remaining).
				remain := len(expected) - offset
				if remain == 0 {
					break
				}
				n := int(expected[(i+2)%len(expected)])%remain + 1
				peek, err := r.Peek(n)
				if err != nil {
					// If we expect enough data but get an error, that's a failure.
					if n <= (len(expected) - offset) {
						t.Errorf("Peek error: %v", err)
					}
				} else {
					if !bytes.Equal(peek, expected[offset:offset+n]) {
						t.Errorf("Peek: got %q, want %q", peek, expected[offset:offset+n])
					}
				}
				// Note: Peek does not advance offset.

			case 3: // ReadBytes
				remain := len(expected) - offset
				if remain == 0 {
					break
				}
				n := int(expected[(i+3)%len(expected)])%remain + 1
				rb, err := r.ReadBytes(n)
				// We allow io.EOF if not enough data remains.
				if err != nil && err != io.EOF {
					t.Errorf("ReadBytes error: %v", err)
				}
				if !bytes.Equal(rb, expected[offset:offset+len(rb)]) {
					t.Errorf("ReadBytes: got %q, want %q", rb, expected[offset:offset+len(rb)])
				}
				offset += len(rb)

			case 4: // Discard
				remain := len(expected) - offset
				if remain == 0 {
					break
				}
				n := int(expected[(i+4)%len(expected)])%remain + 1
				d, err := r.Discard(n)
				if err != nil && err != io.EOF {
					t.Errorf("Discard error: %v", err)
				}
				if d != n {
					t.Errorf("Discard: discarded %d, want %d", d, n)
				}
				offset += d

			case 5: // DiscardUntil
				remain := len(expected) - offset
				if remain == 0 {
					break
				}
				// Choose a target byte from the remaining expected data.
				target := expected[(offset+int(expected[(i+5)%len(expected)]))%len(expected)]
				d, _ := r.DiscardUntil(target)
				// Compute the expected number to discard.
				idx := bytes.IndexByte(expected[offset:], target)
				if idx < 0 {
					idx = remain
				}
				if d != idx {
					t.Errorf("DiscardUntil: discarded %d, want %d", d, idx)
				}
				offset += d
			}
		}
		// Finally, read the rest of the data.
		rest, err := io.ReadAll(r)
		if err != nil && err != io.EOF {
			t.Errorf("Final ReadAll error: %v", err)
		}
		if !bytes.Equal(rest, expected[offset:]) {
			t.Errorf("Final read: got %q, want %q", rest, expected[offset:])
		}
	})
}
