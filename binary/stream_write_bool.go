package binary

func (b *StreamWriter) WriteBool(v bool) error {
	if v {
		return b.WriteByte(1)
	}

	return b.WriteByte(0)
}
