package fast

type StringEncoder interface {
	EncodeString(b *StringBuffer)
}

// Write a type that implements StringEncoder
func (b *StringBuffer) WriteEnc(v StringEncoder) {
	v.EncodeString(b)
}
