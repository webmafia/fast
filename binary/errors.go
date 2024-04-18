package binary

import "errors"

var (
	ErrUnknownValue  = errors.New("unknown value")
	ErrNegativeCount = errors.New("negative count")
)
