package fast

// Write a type that implements StringEncoder
func (b BinaryStream) WriteEnc(v BinaryEncoder) {
	v.BinaryEncode(b)
}
