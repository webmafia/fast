package buffer

import "testing"

func ExampleStringBuffer() {
	buf := NewBuffer(64)

	buf.Str().WritefCb("hello %s", []any{123}, func(b *Buffer, c byte, v any) error {
		b.WriteString("world")
		return nil
	})

	// Output: TODO
}

func BenchmarkStringBuffer(b *testing.B) {
	buf := NewBuffer(64)
	b.ResetTimer()

	for range b.N {
		buf.Str().WritefCb("hello %s", []any{123}, func(b *Buffer, c byte, v any) error {
			b.WriteString("world")
			return nil
		})
		buf.Reset()
	}
}

func BenchmarkStringBufferWritef(b *testing.B) {
	buf := NewBuffer(64)
	b.ResetTimer()

	for range b.N {
		buf.Str().Writef("hello %s", "world")
		buf.Reset()
	}
}
