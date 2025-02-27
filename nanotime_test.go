package fast

import (
	"fmt"
	"testing"
	"time"
)

func ExampleNanotime() {
	start := Nanotime()
	time.Sleep(100 * time.Millisecond)
	end := Nanotime()

	// This will be ~100 ms (most likely 101 ms)
	taken := end - start

	fmt.Println(taken / int64(time.Millisecond))

	// Output: 101
}

func BenchmarkNanotime(b *testing.B) {
	for range b.N {
		_ = Nanotime()
	}
}

func BenchmarkNanotimeSince(b *testing.B) {
	ts := Nanotime()
	b.ResetTimer()

	for range b.N {
		_ = NanotimeSince(ts)
	}
}

func BenchmarkTimeAdd(b *testing.B) {
	ts := time.Now()
	b.ResetTimer()

	for range b.N {
		_ = ts.Add(1 * time.Second)
	}
}
