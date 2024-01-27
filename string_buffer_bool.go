package fast

func (b *StringBuffer) WriteBool(v bool) {
	if v {
		b.buf = append(b.buf, "true"...)
	} else {
		b.buf = append(b.buf, "false"...)
	}
}
