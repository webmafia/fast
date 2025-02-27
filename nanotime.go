package fast

import (
	"time"
	_ "unsafe"
)

// Returns a fast, monotonic time suitable for measuring time taken between two calls.
//
//go:linkname Nanotime runtime.nanotime
func Nanotime() int64

// Returns the duration since a given timestamp acquired from Nanotime.
func NanotimeSince(ts int64) time.Duration {
	return time.Duration(Nanotime() - ts)
}
