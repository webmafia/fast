package buffer

import "errors"

var (
	ErrNegativeCount = errors.New("negative count")
	ErrInvalidValue  = errors.New("invalid value")
	ErrFewArgs       = errors.New("too few arguments")
)
