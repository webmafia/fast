package fast

import (
	"testing"
	"time"
)

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkClockNow(b *testing.B) {
	var c Clock

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Now()
	}
}
