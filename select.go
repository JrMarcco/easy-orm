package orm

import (
	"context"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"strings"
)

type Selector[T any] struct {
	model     *model
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

	var err error
	s.model, err = parseModel(new(T))
	if err != nil {
		return nil, err
	}

	s.sb = &strings.Builder{}
	s.sb.WriteString("SELECT * FROM ")

	if s.tableName == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.tbName)
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

		fd, ok := s.model.fds[exprTyp.name]
		if !ok {
			return errs.InvalidColumnErr(exprTyp.name)
		}

		s.sb.WriteByte('`')
		s.sb.WriteString(fd.colName)
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
			return err
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
			return err
		}

		if _, rok := exprTyp.right.(Predicate); rok {
			s.sb.WriteByte(')')
		}
	default:
		return errs.UnsupportedExprErr
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
