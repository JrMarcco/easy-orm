package errs

import (
	"errors"
	"fmt"
)

var (
	InvalidTypeErr           = errors.New("invalid type, only support struct and first-level pointer")
	InvalidAssignmentErr     = errors.New("invalid assigment type")
	UnsupportedExprErr       = errors.New("unsupported expression type")
	UnsupportedSelectableErr = errors.New("unsupported selectable type")
	EmptyTbNameErr           = errors.New("empty table name")
	EmptyColNameErr          = errors.New("empty column name")
	EmptyInsertRowErr        = errors.New("insert row can not be empty")
	HavingWithoutGroupByErr  = errors.New("having statement can not without group by statement")
)

func InvalidColumnErr(fd string) error {
	return fmt.Errorf("invlaid column '%s'", fd)
}

func InvalidColumnFdErr(fd string) error {
	return fmt.Errorf("invlaid column field '%s'", fd)
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
