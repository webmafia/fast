package strings

import (
	"fmt"
	"time"
	"unsafe"
)

func Example() {
	// Example value
	now := time.Now()

	fmt.Println(getString(unsafe.Pointer(&now)))

	// Output: TODO
}

func getString(ptr unsafe.Pointer) string {
	// Interpret ptr as a pointer to a fmt.Stringer interface.
	// This is only safe if ptr indeed points to a valid fmt.Stringer variable.
	s := *(*fmt.Stringer)(ptr)
	return s.String()
}
