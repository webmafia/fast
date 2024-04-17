package fast

// Write a type that implements StringEncoder
func (b *BinaryBuffer) WriteEnc(v BinaryEncoder) {
	v.BinaryEncode(b)
}
