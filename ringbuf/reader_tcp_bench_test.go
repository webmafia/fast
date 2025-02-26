package ringbuf

import (
	"io"
	"net"
	"strings"
	"testing"
)

const tcpDataSize = 1 * 1024 * 1024 // 1 MB

// generateTCPData returns a string of roughly tcpDataSize bytes.
func generateTCPData() string {
	// Use the alphabet repeated enough times.
	repeats := tcpDataSize / 26
	return strings.Repeat("abcdefghijklmnopqrstuvwxyz", repeats)
}

// persistentTCPServer starts a TCP server on an ephemeral port that accepts one connection
// and then continuously writes the given data in a loop.
func persistentTCPServer(t testing.TB, data []byte) net.Listener {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen on TCP: %v", err)
	}
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		// Continuously write data.
		for {
			_, err := conn.Write(data)
			if err != nil {
				return
			}
		}
	}()
	return ln
}

func BenchmarkRingBufReaderOverTCP(b *testing.B) {
	data := []byte(generateTCPData())
	dataLen := len(data)

	// Start the persistent TCP server.
	ln := persistentTCPServer(b, data)
	defer ln.Close()

	// Dial the connection once, outside the benchmark loop.
	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		b.Fatalf("Dial error: %v", err)
	}
	defer conn.Close()

	b.ResetTimer()

	b.Run("Read", func(b *testing.B) {
		r := NewReader(conn)
		buf := make([]byte, len(data))
		b.SetBytes(int64(dataLen))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// Reset the reader's ring buffer state for a new iteration.
			r.Reset(conn)
			// Read exactly one payload (1 MB).
			_, err := io.ReadFull(r, buf)
			if err != nil {
				b.Fatalf("ReadFull error: %v", err)
			}
		}
	})

	b.Run("ReadByte", func(b *testing.B) {
		r := NewReader(conn)
		b.SetBytes(int64(dataLen))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r.Reset(conn)
			var count int
			for count < dataLen {
				_, err := r.ReadByte()
				if err != nil {
					b.Fatalf("ReadByte error: %v", err)
				}
				count++
			}
			if count != dataLen {
				b.Fatalf("ReadByte: read %d bytes, expected %d", count, dataLen)
			}
		}
	})

	b.Run("ReadBytes_4096", func(b *testing.B) {
		r := NewReader(conn)
		b.SetBytes(int64(dataLen))
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r.Reset(conn)
			var count int
			for count < dataLen {
				_, err := r.ReadBytes(4096)
				if err != nil {
					b.Fatalf("ReadByte error: %v", err)
				}
				count += 4096
			}
			if count < dataLen {
				b.Fatalf("ReadByte: read %d bytes, expected %d", count, dataLen)
			}
		}
	})

}
