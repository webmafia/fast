package fast

type BinaryEncoder interface {
	Encode(b *BinaryBuffer)
}

// Write a type that implements StringEncoder
func (b *BinaryBuffer) WriteEnc(v BinaryEncoder) {
	v.Encode(b)
}
