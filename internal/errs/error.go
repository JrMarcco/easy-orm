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

func InvalidTagContentErr(content string) error {
	return fmt.Errorf("invalid tag content '%s'", content)
}

func EmptyTagKeyErr(content string) error {
	return fmt.Errorf("invalid tag content '%s', key is empty", content)
}

func EmptyTagValErr(content string) error {
	return fmt.Errorf("invalid tag content '%s', val is empty", content)
}
