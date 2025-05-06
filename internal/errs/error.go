package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidModelType = errors.New("[easy-orm] invalid model entity type, only support struct or pointer to struct")
	ErrInvalidTableName = errors.New("[easy-orm] invalid table name")
)

func ErrUnsupportedExpr(expr any) error {
	return fmt.Errorf("[easy-orm] unsupported expression: %v", expr)
}

func ErrInvalidColumn(name string) error {
	return fmt.Errorf("[easy-orm] invalid column name: %s", name)
}
