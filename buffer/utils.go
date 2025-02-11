package buffer

import (
	"math/bits"

	"github.com/webmafia/fast/types"
)

func index(n int) int {
	n--
	n >>= minBitSize

	// Convert n to 0 if n<=0, else n stays n. This ensures idx=0 if n<=0.
	cleanN := n & ^(n >> 31)

	// idx = number of shifts until zero = bits.Len(n)
	idx := bits.Len64(uint64(cleanN))

	// Clamp idx to [0, steps-1]
	m := steps - 1
	mask := (m - idx) >> 31
	idx = (idx & ^mask) | (m & mask)

	return idx
}

func roundPow(n int) int {
	if n <= 1 {
		return 1
	}

	// Start with the number minus one
	n--

	// Spread the highest set bit to the right
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16

	// Add one to get the next power of 2
	return n + 1
}

func sort[T types.Unsigned](a, b []T) {
	for i := 1; i < len(a); i++ {
		j := i - 1

		// Move elements of `a[0...i-1]` that are smaller than `current` (descending order)
		for j >= 0 && a[j] < a[i] {
			a[j+1] = a[j]
			b[j+1] = b[j]
			j--
		}
		// Place the current element at its correct position
		a[j+1] = a[i]
		b[j+1] = b[i]
	}
}
