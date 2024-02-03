package fast

func (b *BinaryBuffer) WriteBool(v bool) {
	if v {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}
