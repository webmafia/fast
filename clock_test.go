package fast

import (
	"context"
	"testing"
	"time"
)

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkClockNow(b *testing.B) {
	c := NewClock(context.Background())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Now()
	}
}
