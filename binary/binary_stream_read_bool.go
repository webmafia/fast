package fast

func (b *BinaryStreamReader) ReadBool() bool {
	var v byte
	v, b.err = b.ReadByte()
	return v != 0
}
