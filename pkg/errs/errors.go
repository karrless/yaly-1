package errs

import (
	"errors"
)

var (
	ErrExpressionNotValid = errors.New("Expression not valid")
	ErrInternal           = errors.New("Internal server error")
)
