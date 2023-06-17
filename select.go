package orm

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	tableName string
	where     []Predicate
	sb        *strings.Builder
	args      []any
}

func (s *Selector[T]) From(tableName string) *Selector[T] {
	s.tableName = tableName
	return s
}

func (s *Selector[T]) Where(predicates ...Predicate) *Selector[T] {
	s.where = predicates
	return s
}

func (s *Selector[T]) Build() (*Statement, error) {

	s.sb = &strings.Builder{}
	s.sb.WriteString("SELECT * FROM ")

	if s.tableName == "" {
		var t T
		typ := reflect.TypeOf(t)

		s.sb.WriteByte('`')
		s.sb.WriteString(typ.Name())
		s.sb.WriteByte('`')
	} else {

		segs := strings.SplitN(s.tableName, ".", 2)

		s.sb.WriteByte('`')
		s.sb.WriteString(segs[0])
		s.sb.WriteByte('`')

		if len(segs) > 1 {
			s.sb.WriteByte('.')
			s.sb.WriteByte('`')
			s.sb.WriteString(segs[1])
			s.sb.WriteByte('`')
		}

	}

	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")

		root := s.where[0]
		for i := 1; i < len(s.where); i++ {
			root = root.And(s.where[i])
		}

		if err := s.buildExpr(root); err != nil {
			return nil, err
		}

	}

	s.sb.WriteByte(';')

	return &Statement{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

// 构建表达式。
// 该过程本是上是一个深度优先遍历二叉树的过程。
func (s *Selector[T]) buildExpr(expr Expression) error {
	if expr == nil {
		return nil
	}

	switch exprTyp := expr.(type) {
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(exprTyp.name)
		s.sb.WriteByte('`')
	case Value:
		s.sb.WriteByte('?')
		s.addArg(exprTyp.val)
	case Predicate:

		if _, lok := exprTyp.left.(Predicate); lok {
			s.sb.WriteByte('(')
		}

		// 递归左子表达式
		if err := s.buildExpr(exprTyp.left); err != nil {
			return nil
		}

		if _, lok := exprTyp.left.(Predicate); lok {
			s.sb.WriteByte(')')
		}

		if exprTyp.left != nil {
			s.sb.WriteByte(' ')
		}
		s.sb.WriteString(string(exprTyp.op))
		if exprTyp.right != nil {
			s.sb.WriteByte(' ')
		}
		if _, rok := exprTyp.right.(Predicate); rok {
			s.sb.WriteByte('(')
		}

		// 递归右子表达式
		if err := s.buildExpr(exprTyp.right); err != nil {
			return nil
		}

		if _, rok := exprTyp.right.(Predicate); rok {
			s.sb.WriteByte(')')
		}
	default:
		return errors.New("unsupported expression type")
	}

	return nil
}

func (s *Selector[T]) addArg(val any) {
	if s.args == nil {
		s.args = make([]any, 0, 4)
	}
	s.args = append(s.args, val)
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
