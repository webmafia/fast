package buffer

import "github.com/webmafia/fast/types"

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
