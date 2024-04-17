package binary

func (b *BufferWriter) WriteBool(v bool) {
	if v {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}
