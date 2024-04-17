package binary

func (b *BufferReader) ReadBool() bool {
	v := b.buf[b.cursor]
	b.cursor++
	return v != 0
}
