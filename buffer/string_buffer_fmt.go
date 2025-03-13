package buffer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/webmafia/fast"
)

func (b StringBuffer) Writef(format string, args ...any) {
	b.B.B = fmt.Appendf(b.B.B, format, args...)
}

func (b StringBuffer) WritefCb(format string, args []any, cb func(b *Buffer, c byte, v any) error) (err error) {
	var cursor int
	var argNum int

	for {
		i := strings.IndexByte(format[cursor:], '%')

		if i < 0 {
			break
		}

		idx := cursor + i
		i = idx + 1
		b.B.WriteString(format[cursor:idx])
		cursor = idx + 2

		// Double % means an escaped %
		if format[i] == '%' {
			b.B.WriteByte('%')
			continue
		}

		if format[i] == '[' {
			end := strings.IndexByte(format[i:], ']')

			if end < 0 {
				return errors.New("missing ']'")
			}

			num, err := strconv.Atoi(format[i+1 : i+end])

			if err != nil {
				return err
			}

			if num < 1 {
				return errors.New("argument number must be at least 1")
			}

			argNum = num
			cursor += end + 1
			i += end + 1
		} else {
			argNum++
		}

		if argNum > len(args) {
			return fmt.Errorf("argument number %d does not exist", argNum)
		}

		c := format[i]

		if argNum > len(args) {
			return ErrFewArgs
		}

		if err = cb(b.B, c, *fast.NoescapeVal(&args[argNum-1])); err != nil {
			return
		}
	}

	b.B.WriteString(format[cursor:])
	return
}
