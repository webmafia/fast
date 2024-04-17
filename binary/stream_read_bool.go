package binary

func (b *StreamReader) ReadBool() bool {
	var v byte
	v, b.err = b.ReadByte()
	return v != 0
}
