package errs

import (
	"errors"
	"fmt"
)

var (
	InvalidTypeErr     = errors.New("invalid type, only support struct and first-level pointer")
	UnsupportedExprErr = errors.New("unsupported expression type")
)

func InvalidColumnErr(col string) error {
	return fmt.Errorf("invlaid column '%s'", col)
}
