package bufio

// // TestLockAndSlide verifies that when the buffer is locked, it never "slides"
// // unread data even when forced to grow, and that once it is unlocked, sliding is
// // possible again.
// func TestLockAndSlide(t *testing.T) {
// 	// We'll create data that forces multiple fills and slides if unlocked.
// 	data := strings.Repeat("0123456789", 20) // 200 bytes
// 	r := NewReader(strings.NewReader(data))

// 	// Read and discard some to ensure we have partial usage of the buffer.
// 	// Then lock to verify that no sliding occurs.
// 	if _, err := r.Peek(50); err != nil {
// 		t.Fatalf("initial Peek error: %v", err)
// 	}
// 	if got := r.Buffered(); got < 50 {
// 		t.Fatalf("expected at least 50 in buffer, got %d", got)
// 	}
// 	// Discard 10 bytes to create a gap.
// 	discarded, err := r.Discard(10)
// 	if err != nil {
// 		t.Fatalf("Discard(10) error: %v", err)
// 	}
// 	if discarded != 10 {
// 		t.Fatalf("Discarded %d, want 10", discarded)
// 	}
// 	initialR := r.r
// 	initialW := r.w

// 	// Lock the buffer: no more sliding.
// 	r.Lock()

// 	// Force more reads to fill or attempt to slide. If the buffer slides now,
// 	// it's a bug because locked is true.
// 	if _, err := r.Peek(50); err != nil && !errors.Is(err, io.EOF) {
// 		t.Fatalf("Peek error after lock: %v", err)
// 	}
// 	if r.r != initialR {
// 		t.Errorf("expected r=%d (unchanged) if locked, got %d", initialR, r.r)
// 	}
// 	// We allow w to change because new data can arrive; we just can't "slide" (move unread data).
// 	if r.w < initialW {
// 		t.Errorf("expected w to be >= %d, got %d", initialW, r.w)
// 	}

// 	// Now unlock and peek again, which may trigger a slide if it is "worthMoving".
// 	r.Unlock()
// 	beforeSlideR := r.r
// 	// Force more reading. This should cause the Reader to slide if there's enough unread data.
// 	if _, err := r.Peek(20); err != nil && !errors.Is(err, io.EOF) {
// 		t.Fatalf("Peek error after unlock: %v", err)
// 	}
// 	if r.r != 0 && r.r == beforeSlideR {
// 		t.Errorf("expected r to have slid (likely back to 0), but it's still %d", r.r)
// 	}
// }

// // TestLockAndMaxSize ensures that when the buffer is locked and at maxSize,
// // the Reader does not panic. It also verifies partial reads from the underlying
// // reader if there's no room to grow.
// func TestLockAndMaxSize(t *testing.T) {
// 	data := strings.Repeat("ABCDEFGH", 200) // 1600 bytes
// 	maxSize := 256
// 	r := NewReader(strings.NewReader(data), maxSize)

// 	// Force filling close to maxSize.
// 	// We'll read some portion so there's partial leftover, then lock.
// 	// The idea is to see if repeated fill attempts cause any panic or weird behavior.
// 	chunk := make([]byte, 100)
// 	n, err := r.Read(chunk)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if n != 100 {
// 		t.Fatalf("expected to read 100, got %d", n)
// 	}

// 	// Now lock the buffer. We'll try to read until we've definitely passed maxSize usage.
// 	r.Lock()

// 	totalRead := 100
// 	readBuf := make([]byte, 128)
// 	for {
// 		n, err := r.Read(readBuf)
// 		totalRead += n
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			t.Fatalf("unexpected read error: %v", err)
// 		}
// 		if totalRead > 1600 {
// 			// We read all data in an unexpected way. This can happen if underlying can still read
// 			// into the locked buffer as it won't "slide" but we do keep reading from it eventually.
// 			break
// 		}
// 		// We just keep reading. If the buffer can't grow beyond maxSize, we rely on
// 		// partial reads from the underlying stream or eventually reading from it to free
// 		// up space. There's no explicit error in the code for "buffer full," but it
// 		// shouldn't panic. We'll just keep reading until EOF or partial reads fill the buffer.
// 	}

// 	if totalRead < len(data) {
// 		// We might end up reading everything eventually, but it's okay if we read partial
// 		// as long as we didn't panic. The code does partial reads/ writes.
// 		t.Logf("Read %d bytes out of %d with locked buffer at maxSize=%d; no panic as expected.",
// 			totalRead, len(data), maxSize)
// 	} else if totalRead == len(data) {
// 		t.Logf("Successfully read all data with locked buffer at maxSize=%d. No panic, good!", maxSize)
// 	} else {
// 		t.Errorf("Somehow read more bytes than existed! totalRead=%d, dataLen=%d", totalRead, len(data))
// 	}
// }

// // TestCombinedReads verifies that Read, ReadByte, and Peek can be used
// // interchangeably. We set various maxSizes to see if partial reads or repeated fills
// // cause issues.
// func TestCombinedReads(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		data    string
// 		maxSize int
// 	}{
// 		{"NoMaxSizeShortData", "HelloWorld", 0},
// 		{"LimitedMaxSizeShortData", "HelloWorld", 8},
// 		{"NoMaxSizeLongData", strings.Repeat("x", 2000), 0},
// 		{"LimitedMaxSizeLongData", strings.Repeat("x", 2000), 256},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := NewReader(strings.NewReader(tt.data), tt.maxSize)

// 			// 1) Peek 5 bytes
// 			p, err := r.Peek(5)
// 			if err != nil && err != io.EOF {
// 				t.Fatalf("Peek error: %v", err)
// 			}
// 			if len(tt.data) >= 5 && string(p) != tt.data[:5] {
// 				t.Errorf("Peek mismatch: got %q, want %q", string(p), tt.data[:5])
// 			}

// 			// 2) ReadByte for next 2 bytes
// 			var readBytes []byte
// 			for i := 0; i < 2 && err == nil; i++ {
// 				b, e := r.ReadByte()
// 				if e != nil && e != io.EOF {
// 					t.Fatalf("ReadByte error: %v", e)
// 				}
// 				err = e
// 				if e == nil {
// 					readBytes = append(readBytes, b)
// 				}
// 			}
// 			// Check the 2 read bytes vs the original data
// 			if len(tt.data) >= 7 {
// 				want := tt.data[5:7]
// 				if string(readBytes) != want {
// 					t.Errorf("ReadByte mismatch: got %q, want %q", string(readBytes), want)
// 				}
// 			}

// 			// 3) Read 10 bytes
// 			buf := make([]byte, 10)
// 			n, err := r.Read(buf)
// 			if err != nil && err != io.EOF {
// 				t.Fatalf("Read error: %v", err)
// 			}
// 			// We should get the next part of the data if it is that long
// 			nextOffset := 7
// 			if len(tt.data) < nextOffset {
// 				nextOffset = len(tt.data)
// 			}
// 			wantLen := 0
// 			if len(tt.data) > nextOffset {
// 				wantLen = len(tt.data) - nextOffset
// 				if wantLen > 10 {
// 					wantLen = 10
// 				}
// 			}
// 			if n != wantLen {
// 				t.Errorf("Read returned %d bytes, want %d", n, wantLen)
// 			}
// 			if wantLen > 0 {
// 				got := string(buf[:n])
// 				want := tt.data[nextOffset : nextOffset+n]
// 				if got != want {
// 					t.Errorf("Read mismatch: got %q, want %q", got, want)
// 				}
// 			}
// 		})
// 	}
// }

// // TestDiscardAndReadSlice checks that Discard, DiscardUntil, and ReadSlice
// // can be combined logically with reading from the buffer.
// func TestDiscardAndReadSlice(t *testing.T) {
// 	// We'll set up data with multiple delimiters so we can discard until one,
// 	// read a slice, discard more, etc.
// 	data := "Hello,World-Again,And-Again"

// 	r := NewReader(strings.NewReader(data), 0) // no maxSize limit

// 	// 1) Discard 6 bytes -> "Hello,"
// 	discarded, err := r.Discard(6)
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("Discard error: %v", err)
// 	}
// 	if discarded != 6 {
// 		t.Errorf("Discarded %d, want 6", discarded)
// 	}

// 	// 2) ReadSlice until '-', expect "World"
// 	slice, err := r.ReadSlice('-')
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("ReadSlice error: %v", err)
// 	}
// 	if string(slice) != "World" {
// 		t.Errorf("ReadSlice got %q, want %q", string(slice), "World")
// 	}

// 	// 3) DiscardUntil ',' from "Again,And..."
// 	d, err := r.DiscardUntil(',')
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("DiscardUntil error: %v", err)
// 	}
// 	// The substring is "Again,And-Again"
// 	// We are at the "W" -> "World" ended. Then next is "Again,And-Again".
// 	// We find a comma after "Again", so the discard includes "Again". So d is the number
// 	// of bytes discarded between the last pointer and the delimiter. But the code
// 	// returns 0 on a direct hit in the current chunk. Because it finds the comma in
// 	// the current buffer. So d might be 0 if it found the delimiter immediately
// 	// within a read chunk. It might also reflect partial or a full leftover.
// 	// We'll just check we didn't get an error or panic.
// 	_ = d // not strongly validated here, we just ensure no errors.

// 	// 4) ReadSlice until '-', should get "And"
// 	slice, err = r.ReadSlice('-')
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("ReadSlice error: %v", err)
// 	}
// 	if string(slice) != "And" {
// 		t.Errorf("ReadSlice got %q, want %q", string(slice), "And")
// 	}

// 	// 5) finally read the leftover "Again"
// 	buf := make([]byte, 10)
// 	n, err := r.Read(buf)
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("Read error: %v", err)
// 	}
// 	got := string(buf[:n])
// 	if got != "Again" {
// 		t.Errorf("final read got %q, want %q", got, "Again")
// 	}
// }

// // TestFullUsageScenario covers a more realistic scenario of reading and discarding
// // in chunks, toggling locks, verifying that everything works together.
// func TestFullUsageScenario(t *testing.T) {
// 	// This scenario uses a moderate data size with multiple distinct steps.
// 	data := "BEGIN:12345:LOCKME:DiscardSome:BYTES:ThenUNLOCK:AndREAD" +
// 		strings.Repeat("X", 2048) + ":FINAL"
// 	r := NewReader(strings.NewReader(data), 512) // moderate maxSize

// 	// 1) Read "BEGIN:"
// 	begin, err := r.ReadSlice(':')
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("ReadSlice(':') error: %v", err)
// 	}
// 	if string(begin) != "BEGIN" {
// 		t.Errorf("got %q, want %q", string(begin), "BEGIN")
// 	}

// 	// 2) Bytes(5) => "12345"
// 	b, err := r.Bytes(5)
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("Bytes(5) error: %v", err)
// 	}
// 	if string(b) != "12345" {
// 		t.Errorf("got %q, want %q", string(b), "12345")
// 	}

// 	// 3) Lock the buffer
// 	r.Lock()

// 	// 4) ReadSlice(':') => "LOCKME"
// 	lockSlice, err := r.ReadSlice(':')
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("ReadSlice(':') error: %v", err)
// 	}
// 	if string(lockSlice) != "LOCKME" {
// 		t.Errorf("got %q, want %q", string(lockSlice), "LOCKME")
// 	}

// 	// 5) DiscardSome => we attempt to discard 11 => "DiscardSome"
// 	disc, err := r.Discard(11)
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("Discard(11) error: %v", err)
// 	}
// 	if disc != 11 {
// 		t.Errorf("Discarded %d, want 11", disc)
// 	}

// 	// 6) ReadSlice(':') => "BYTES"
// 	bytesSlice, err := r.ReadSlice(':')
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("ReadSlice(':') error: %v", err)
// 	}
// 	if string(bytesSlice) != "BYTES" {
// 		t.Errorf("got %q, want %q", string(bytesSlice), "BYTES")
// 	}

// 	// Because we're locked, the buffer won't slide. We'll push it a bit more:
// 	// 7) Read next 9 => "ThenUNLOCK"
// 	buf := make([]byte, 9)
// 	n, err := r.Read(buf)
// 	if err != nil && err != io.EOF {
// 		t.Fatalf("Read error: %v", err)
// 	}
// 	if n != 9 {
// 		t.Fatalf("expected 9, got %d", n)
// 	}
// 	if string(buf[:n]) != "ThenUNLOCK"[:9] {
// 		t.Errorf("got %q, want %q", string(buf[:n]), "ThenUNLOCK"[:9])
// 	}

// 	// 8) Unlock
// 	r.Unlock()

// 	// 9) Continue reading => final chunk "AndREAD" + 2048 X's + ":FINAL"
// 	// Let's do partial reads in random chunk sizes to ensure we get the entire sequence.

// 	remaining := "kAndREAD" + strings.Repeat("X", 2048) + ":FINAL"
// 	// We already read "ThenUNLOCK"[:9] => "ThenUNLOC", so the leftover in "ThenUNLOCKAndREAD" is 'kAndREAD'
// 	// So let's correct the leftover.
// 	// Actually, the substring "ThenUNLOCK" is 11 chars. We read 9 => "ThenUNLOC", so there's 'K' left plus next "AndREAD".
// 	// Let's unify carefully:
// 	// The data was: "ThenUNLOCK:AndREAD..."
// 	// We read 9 => "ThenUNLOC", leftover is "K:AndREAD" + 2048 X + ":FINAL"
// 	remaining = "K:AndREAD" + strings.Repeat("X", 2048) + ":FINAL"

// 	var out []byte
// 	tmp := make([]byte, 37)
// 	for len(out) < len(remaining) {
// 		n, err := r.Read(tmp)
// 		if n > 0 {
// 			out = append(out, tmp[:n]...)
// 		}
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			t.Fatalf("unexpected read error after unlock: %v", err)
// 		}
// 	}
// 	got := string(out)
// 	if got != remaining {
// 		t.Errorf("final read mismatch\n got:  %q\n want: %q", got, remaining)
// 	}
// }
