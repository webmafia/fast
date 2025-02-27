package ringbuf

import (
	"bytes"
	"io"
	"testing"
)

func FuzzReadBytes(f *testing.F) {
	// Seed the fuzzer with some initial test cases.
	f.Add("hello world")
	f.Add("The quick brown fox jumps over the lazy dog")
	f.Add("1234567890")

	f.Fuzz(func(t *testing.T, input string) {
		// Truncate input to BufferSize (4096) if necessary.
		if len(input) > BufferSize {
			input = input[:BufferSize]
		}

		// Create a new RingBuf and reset it.
		var rb RingBuf
		rb.Reset()

		// Write the input data to the ring buffer.
		data := []byte(input)
		n, err := rb.Write(data)
		if err != nil {
			t.Fatalf("Write error: %v", err)
		}
		if n != len(data) {
			t.Fatalf("expected to write %d bytes, wrote %d", len(data), n)
		}

		// Read the bytes back in chunks of variable sizes.
		var result []byte
		expected := data
		remaining := len(expected)

		// To vary the chunk size, we'll use a simple loop where the chunk size
		// is determined by the remaining bytes. For example, we use modulo arithmetic.
		for remaining > 0 {
			// Choose a chunk size between 1 and 10 (or less if remaining is smaller).
			chunkSize := (remaining % 10) + 1
			if chunkSize > remaining {
				chunkSize = remaining
			}

			// Read the next chunk.
			chunk, err := rb.ReadBytes(chunkSize)
			if err != nil && err != io.EOF {
				t.Fatalf("ReadBytes error: %v", err)
			}

			// The length of the returned chunk should exactly match the requested length,
			// unless we have reached the end of the data.
			if len(chunk) != chunkSize {
				t.Fatalf("expected chunk of length %d, got %d", chunkSize, len(chunk))
			}

			result = append(result, chunk...)
			remaining -= len(chunk)
		}

		// Compare the complete output with the input data.
		if !bytes.Equal(result, expected) {
			t.Fatalf("expected output %q, got %q", expected, result)
		}
	})
}

// padToFullBuffer ensures the returned slice is exactly BufferSize bytes.
// If src is shorter, it repeats src until BufferSize is reached.
func padToFullBuffer(src []byte) []byte {
	buf := make([]byte, BufferSize)
	for i := 0; i < int(BufferSize); i++ {
		buf[i] = src[i%len(src)]
	}
	return buf
}

func FuzzReadBytesWrapSlack(f *testing.F) {
	// Seed fuzzer with some initial test cases.
	f.Add([]byte("abcdefghijklmnopqrstuvwxyz"), []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), 123, 2000)
	f.Add([]byte("0123456789"), []byte("9876543210"), 100, 3500)
	f.Add([]byte("The quick brown fox jumps over the lazy dog"), []byte("!@#$%^&*()"), 2048, 3000)

	f.Fuzz(func(t *testing.T, firstData, secondData []byte, readOffset, readLength int) {
		// For a valid wrap-around test, both firstData and secondData must be non-empty.
		if len(firstData) == 0 || len(secondData) == 0 {
			return
		}
		// Ensure firstData fills the logical ring exactly.
		firstData = padToFullBuffer(firstData)

		// Clamp readOffset to be within 1 and BufferSize-1.
		if readOffset <= 0 {
			readOffset = 1
		} else if readOffset >= BufferSize {
			readOffset = BufferSize - 1
		}

		// After writing firstData, the unread region is the full 4096 bytes.
		// After reading readOffset bytes, the unread region becomes firstData[readOffset:].
		// That leaves free space equal to readOffset bytes (at the beginning of the ring).
		freeSpace := readOffset

		// Clamp secondData to freeSpace.
		if len(secondData) > freeSpace {
			secondData = secondData[:freeSpace]
		}

		// The unread region now is:
		//   - Tail: firstData[readOffset:] (length = BufferSize - readOffset)
		//   - Head: secondData (length = len(secondData))
		// Total unread bytes:
		totalUnread := (BufferSize - readOffset) + len(secondData)
		if totalUnread == 0 {
			return
		}

		// To force the slack branch in Peek, we need to read more than the contiguous tail.
		tailLen := BufferSize - readOffset
		if readLength <= tailLen {
			readLength = tailLen + 1
		}
		if readLength > totalUnread {
			readLength = totalUnread
		}

		// Initialize the ring buffer.
		var rb RingBuf
		rb.Reset()

		// Write the firstData into the ring buffer.
		n, err := rb.Write(firstData)
		if err != nil || n != len(firstData) {
			t.Fatalf("failed writing firstData: err=%v, written=%d, expected=%d", err, n, len(firstData))
		}

		// Advance the read pointer by reading 'readOffset' bytes.
		consumed, err := rb.ReadBytes(readOffset)
		if err != nil || len(consumed) != readOffset {
			t.Fatalf("failed reading readOffset (%d bytes): err=%v, got %d bytes", readOffset, err, len(consumed))
		}

		// Write secondData into the ring buffer.
		n, err = rb.Write(secondData)
		if err != nil || n != len(secondData) {
			t.Fatalf("failed writing secondData: err=%v, written=%d, expected=%d", err, n, len(secondData))
		}

		// Build the expected unread data:
		//   - Tail of firstData, from readOffset to end.
		//   - Then secondData.
		expected := append(append([]byte(nil), firstData[readOffset:]...), secondData...)

		// Now, call ReadBytes to retrieve readLength bytes.
		result, err := rb.ReadBytes(readLength)
		if err != nil {
			t.Fatalf("failed reading readLength (%d bytes): err=%v", readLength, err)
		}
		if len(result) != readLength {
			t.Fatalf("expected %d bytes, got %d", readLength, len(result))
		}

		// Validate that the returned bytes match the expected bytes.
		if !bytes.Equal(result, expected[:readLength]) {
			t.Fatalf("data mismatch:\nexpected: %v\ngot:      %v", expected[:readLength], result)
		}
	})
}
