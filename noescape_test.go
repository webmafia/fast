package fast_test

import (
	"testing"
	"unsafe"

	"github.com/webmafia/fast"
)

func TestNoescapeUnsafe(t *testing.T) {
	x := 42
	ptr := unsafe.Pointer(&x)
	result := fast.NoescapeUnsafe(ptr)

	if result != ptr {
		t.Errorf("NoescapeUnsafe did not return the same pointer: got %v, want %v", result, ptr)
	}
}

func TestNoescape(t *testing.T) {
	x := 123
	result := fast.Noescape(&x)

	if result != &x {
		t.Errorf("Noescape did not return the same value: got %d, want %d", result, x)
	}
}

type dummy struct {
	a, b, c int
}

var sink *dummy

func BenchmarkEscape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := dummy{1, 2, 3}
		sink = &x // Escapes, thus forces heap allocation
	}
}

func BenchmarkNoescape(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := dummy{1, 2, 3}
		sink = fast.Noescape(&x) // Hides from escape analysis
	}
}
