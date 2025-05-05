package easyorm

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

var (
	errInvalidTableName = errors.New("invalid table name")
)

// selectable marker interface, used to identify optional query columns( e.g., columns, aggregate functions ).
type selectable interface {
	selectable()
}

var _ Querier[any] = (*Selector[any])(nil)

type Selector[T any] struct {
	tableName string
	where     []Predicate

	sqlBuilder strings.Builder
	args       []any
}

func (s *Selector[T]) FindOne(ctx context.Context) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) FindMulti(ctx context.Context) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) From(tableName string) *Selector[T] {
	s.tableName = tableName
	return s
}

func (s *Selector[T]) Where(pds ...Predicate) *Selector[T] {
	s.where = append(s.where, pds...)
	return s
}

func (s *Selector[T]) Build() (*Statement, error) {
	s.sqlBuilder.WriteString("SELECT * FROM ")

	if s.tableName == "" {
		var t T
		typeOfT := reflect.TypeOf(t)
		s.sqlBuilder.WriteByte('`')
		s.sqlBuilder.WriteString(typeOfT.Name())
		s.sqlBuilder.WriteByte('`')
	} else {
		segments := strings.Split(s.tableName, ".")
		if len(segments) > 2 {
			return nil, errInvalidTableName
		}

		s.sqlBuilder.WriteByte('`')
		s.sqlBuilder.WriteString(segments[0])
		s.sqlBuilder.WriteByte('`')

		if len(segments) == 2 {
			s.sqlBuilder.WriteByte('.')
			s.sqlBuilder.WriteByte('`')
			s.sqlBuilder.WriteString(segments[1])
			s.sqlBuilder.WriteByte('`')
		}
	}

	lenOfWhere := len(s.where)
	if lenOfWhere > 0 {
		s.sqlBuilder.WriteString(" WHERE ")

		p := s.where[0]
		for i := 1; i < lenOfWhere; i++ {
			p = p.And(s.where[i])
		}
		if err := s.buildExpr(p); err != nil {
			return nil, err
		}
	}

	s.sqlBuilder.WriteByte(';')
	return &Statement{
		Sql:  s.sqlBuilder.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpr(expr Expression) error {
	if expr == nil {
		return nil
	}

	switch exprType := expr.(type) {
	case Predicate:
		if _, ok := exprType.left.(Predicate); ok {
			s.sqlBuilder.WriteByte('(')
		}
		if err := s.buildExpr(exprType.left); err != nil {
			return err
		}
		if _, ok := exprType.left.(Predicate); ok {
			s.sqlBuilder.WriteByte(')')
		}

		if exprType.op != "" {
			if exprType.left != nil {
				s.sqlBuilder.WriteByte(' ')
			}
			s.sqlBuilder.WriteString(exprType.op.String())
			if exprType.right != nil {
				s.sqlBuilder.WriteByte(' ')
			}

			if _, ok := exprType.right.(Predicate); ok {
				s.sqlBuilder.WriteByte('(')
			}
			if err := s.buildExpr(exprType.right); err != nil {
				return err
			}
			if _, ok := exprType.right.(Predicate); ok {
				s.sqlBuilder.WriteByte(')')
			}
		}
	case Column:
		s.sqlBuilder.WriteByte('`')
		s.sqlBuilder.WriteString(exprType.name)
		s.sqlBuilder.WriteByte('`')
	case ColumnValue:
		s.sqlBuilder.WriteByte('?')
		s.addArgs(exprType.value)
	}
	return nil
}

func (s *Selector[T]) addArgs(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
}

func NewSelector[T any]() *Selector[T] {
	return &Selector[T]{}
}
