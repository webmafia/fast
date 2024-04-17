package binary

func (b StreamWriter) WriteBool(v bool) {
	if v {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}
}