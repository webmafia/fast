package fast

func NewStringBufferPool() *Pool[StringBuffer] {
	return NewPool[StringBuffer](nil, func(b *StringBuffer) {
		b.Reset()
	})
}
