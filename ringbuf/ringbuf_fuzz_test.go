// file: ringbuf_fuzz_no_rewind_test.go
package ringbuf

import (
	"bytes"
	"math/rand"
	"testing"
)

// Operation types (no rewind).
const (
	opWrite = iota
	opRead
	opReadBytes
	opToggleManual
	opFlush
)

// FuzzRingBufNoRewind exercises the ring buffer without ever calling Rewind.
func FuzzRingBufNoRewind(f *testing.F) {
	// Add some seed inputs.
	f.Add([]byte{})                                  // empty
	f.Add([]byte("hello"))                           // small
	f.Add([]byte("some random input for ring"))      // medium
	f.Add(bytes.Repeat([]byte("X"), 8192))           // large
	f.Add(bytes.Repeat([]byte("1234567890"), 10000)) // more than BufferSize

	f.Fuzz(func(t *testing.T, data []byte) {
		rng := rand.New(rand.NewSource(int64(len(data)) * 12345))
		rb := &RingBuf{}
		rb.Reset()

		// Our oracle to compare reads against.
		var oracle bytes.Buffer

		// Save slices returned by ReadBytes to confirm they remain intact
		// until a Flush occurs.
		var savedRefs []struct {
			dataRef  []byte
			snapshot []byte
		}

		manualOn := false
		N := 300 // number of operations

		for i := 0; i < N; i++ {
			op := rng.Intn(5) // only 0..4 now
			switch op {
			case opWrite:
				chunkSize := rng.Intn(200) + 1
				randomChunk := make([]byte, chunkSize)
				for j := 0; j < chunkSize; j++ {
					randomChunk[j] = byte(rng.Intn(256))
				}

				n, err := rb.Write(randomChunk)
				if n > 0 {
					oracle.Write(randomChunk[:n])
				}

				// Check no overwrite of start.
				if rb.write-rb.start > BufferSize {
					t.Fatalf("overwrite detected: write=%d start=%d (cap=%d)",
						rb.write, rb.start, BufferSize)
				}
				// Acceptable errors: nil, io.ErrShortBuffer, or EOF.
				if err != nil && err.Error() != "short buffer" && err.Error() != "EOF" {
					t.Fatalf("unexpected write error: %v", err)
				}

			case opRead:
				beforeStart := rb.start
				chunkSize := rng.Intn(200) + 1
				p := make([]byte, chunkSize)
				n, err := rb.Read(p)
				afterStart := rb.start
				p = p[:n]

				if err != nil && err.Error() != "EOF" {
					t.Fatalf("unexpected read error: %v", err)
				}

				// Compare with oracle
				want := oracle.Next(n)
				if !bytes.Equal(p, want) {
					t.Fatalf("mismatch in read vs oracle:\n got=%v\nwant=%v", p, want)
				}

				// If manual flush is on, confirm `start` did NOT move
				// (it should only move on Flush).
				if manualOn && afterStart != beforeStart {
					t.Fatalf("manual flush on, but start moved automatically!\n before=%d after=%d",
						beforeStart, afterStart)
				}

			case opReadBytes:
				beforeStart := rb.start
				chunkSize := rng.Intn(200) + 1
				rbData, err := rb.ReadBytes(chunkSize)
				afterStart := rb.start

				if err != nil && err.Error() != "EOF" {
					t.Fatalf("unexpected error from ReadBytes: %v", err)
				}

				if err == nil {
					// Compare with oracle
					want := oracle.Next(len(rbData))
					if !bytes.Equal(rbData, want) {
						t.Fatalf("readbytes mismatch vs oracle:\n got=%v\nwant=%v", rbData, want)
					}

					// Save a copy for verifying immutability until Flush
					copySlice := append([]byte(nil), rbData...)
					savedRefs = append(savedRefs, struct {
						dataRef  []byte
						snapshot []byte
					}{rbData, copySlice})

					// If manual flush is on, confirm `start` did NOT move
					if manualOn && afterStart != beforeStart {
						t.Fatalf("manual flush on, but start moved automatically on ReadBytes!\n before=%d after=%d",
							beforeStart, afterStart)
					}
				}

			case opToggleManual:
				manualOn = !manualOn
				rb.SetManualFlush(manualOn)

			case opFlush:
				rb.Flush()
				// All old references can become invalid now.
				savedRefs = nil
			}

			// Check that data from ReadBytes hasn't changed unless we've flushed.
			for _, ref := range savedRefs {
				if !bytes.Equal(ref.dataRef, ref.snapshot) {
					t.Fatalf("ReadBytes data changed unexpectedly!\n got=%v\nwant=%v",
						ref.dataRef, ref.snapshot)
				}
			}
		}
	})
}
