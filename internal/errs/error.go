package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidType           = errors.New("invalid type, only support struct and first-level pointer")
	ErrInvalidAssignment     = errors.New("invalid assigment type")
	ErrUnsupportedExpr       = errors.New("unsupported expression type")
	ErrUnsupportedSelectable = errors.New("unsupported selectable type")
	ErrEmptyTbName           = errors.New("empty table name")
	ErrEmptyColName          = errors.New("empty column name")
	ErrEmptyInsertRow        = errors.New("insert row can not be empty")
	ErrNoEligibleRows        = errors.New("no eligible rows")
	ErrHavingWithoutGroupBy  = errors.New("having statement can not without group by statement")
	ErrUnsafeDelete          = errors.New("unsafe operation, delete without where")
)

func ErrInvalidColumn(fd string) error {
	return fmt.Errorf("invlaid column '%s'", fd)
}

func ErrInvalidColumnFd(fd string) error {
	return fmt.Errorf("invlaid column field '%s'", fd)
}

func ErrInvalidTagContent(content string) error {
	return fmt.Errorf("invalid tag content '%s'", content)
}

func ErrInvalidTbRefType(tbRef any) error {
	return fmt.Errorf("invalid table reference type: %v", tbRef)
}

func ErrEmptyTagKey(content string) error {
	return fmt.Errorf("invalid tag content '%s', key is empty", content)
}

func ErrEmptyTagVal(content string) error {
	return fmt.Errorf("invalid tag content '%s', val is empty", content)
}

func ErrRollback(bizErr error, rollbackErr error, bizPanicked bool) error {
	return fmt.Errorf(
		"fial to rollback, business error: %w, rollback error: %s, business panicked: %t",
		bizErr, rollbackErr, bizPanicked,
	)
}
