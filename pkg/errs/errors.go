package errs

import (
	"errors"
)

var (
	ErrExpressionNotValid = errors.New("Expression not valid")  // 422
	ErrInternal           = errors.New("Internal server error") // 500
	ErrDivisionByZero     = errors.New("Division by zero")
)
