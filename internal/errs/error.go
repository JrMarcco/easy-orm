package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidModelType = errors.New("[easy-orm] invalid model entity type, only support struct or pointer to struct")
	ErrEligibleRow      = errors.New("[easy-orm] eligible row not found")
)

func ErrUnsupportedExpr(expr any) error {
	return fmt.Errorf("[easy-orm] unsupported expression: %v", expr)
}

func ErrInvalidField(fieldName string) error {
	return fmt.Errorf("[easy-orm] invalid field: %s", fieldName)
}

func ErrInvalidTable(name string) error {
	return fmt.Errorf("[easy-orm] invalid table: %s", name)
}

func ErrInvalidColumn(name string) error {
	return fmt.Errorf("[easy-orm] invalid column: %s", name)
}

func ErrInvalidTag(tagPair string) error {
	return fmt.Errorf("[easy-orm] invalid tag: %s", tagPair)
}

func ErrRollback(bizErr, rbErr error, bizPanicked bool) error {
	return fmt.Errorf(
		"[easy-orm] failed to rollback for biz error: %v, rollback error: %v, business panicked: %v",
		bizErr, rbErr, bizPanicked,
	)
}
