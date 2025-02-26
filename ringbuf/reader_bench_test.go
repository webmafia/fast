package ringbuf

import (
	"bufio"
	"io"
	"strings"
	"testing"
)

const benchDataSize = 1 * 1024 * 1024 // 1 MB

// generateData returns a string of roughly benchDataSize bytes.
func generateData() string {
	// 26 letters repeated enough times.
	repeats := benchDataSize / 26
	return strings.Repeat("abcdefghijklmnopqrstuvwxyz", repeats)
}

//
// Benchmarks for our custom ringbuf Reader
//

func BenchmarkRingBufReaderRead(b *testing.B) {
	data := generateData()
	// Set bytes processed per iteration.
	b.SetBytes(int64(len(data)))
	buf := make([]byte, 4096)
	sr := strings.NewReader(data)
	r := NewReader(sr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sr.Reset(data)
		r.Reset(r)
		for {
			n, err := r.Read(buf)
			if err == io.EOF {
				break
			}
			if n == 0 {
				break
			}
		}
	}
}

func BenchmarkRingBufReaderReadBytes(b *testing.B) {
	data := generateData()
	b.SetBytes(int64(len(data)))
	sr := strings.NewReader(data)
	r := NewReader(sr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sr.Reset(data)
		r.Reset(r)
		for {
			_, err := r.ReadBytes(4096)
			if err == io.EOF {
				break
			}
		}
	}
}

func BenchmarkRingBufReaderReadByte(b *testing.B) {
	data := generateData()
	b.SetBytes(int64(len(data)))
	sr := strings.NewReader(data)
	r := NewReader(sr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sr.Reset(data)
		r.Reset(r)
		for {
			_, err := r.ReadByte()
			if err == io.EOF {
				break
			}
		}
	}
}

//
// Benchmarks for bufio.Reader (the standard library)
//

func BenchmarkBufioReaderRead(b *testing.B) {
	data := generateData()
	b.SetBytes(int64(len(data)))
	buf := make([]byte, 4096)
	sr := strings.NewReader(data)
	// Using bufio.NewReader for the standard benchmark.
	r := bufio.NewReader(sr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sr.Reset(data)
		r.Reset(sr)
		for {
			n, err := r.Read(buf)
			if err == io.EOF {
				break
			}
			if n == 0 {
				break
			}
		}
	}
}

func BenchmarkBufioReaderReadByte(b *testing.B) {
	data := generateData()
	b.SetBytes(int64(len(data)))
	sr := strings.NewReader(data)
	r := bufio.NewReader(sr)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sr.Reset(data)
		r.Reset(sr)
		for {
			_, err := r.ReadByte()
			if err == io.EOF {
				break
			}
		}
	}
}
