package buffer

type StringBuffer struct {
	B *Buffer
}

func (b StringBuffer) WriteBool(v bool) {
	if v {
		b.B.B = append(b.B.B, "true"...)
	} else {
		b.B.B = append(b.B.B, "false"...)
	}
}

func (b StringBuffer) WriteString(str string) {
	b.B.WriteString(str)
}
