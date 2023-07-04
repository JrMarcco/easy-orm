package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/model"
	"strings"
)

type builder struct {
	tbName string
	model  *model.Model
	sb     *strings.Builder
	args   []any
}

func newBuilder() *builder {
	return &builder{
		sb: &strings.Builder{},
	}
}

// 构建表达式。
// 该过程本是上是一个深度优先遍历二叉树的过程。
func (b *builder) buildExpr(expr Expression) error {
	if expr == nil {
		return nil
	}

	switch exprTyp := expr.(type) {
	case Column:

		if err := b.buildCol(exprTyp); err != nil {
			return err
		}
		//fd, ok := b.model.Fds[exprTyp.name]
		//if !ok {
		//	return errs.InvalidColumnFdErr(exprTyp.name)
		//}
		//
		//b.sb.WriteByte('`')
		//b.sb.WriteString(fd.ColName)
		//b.sb.WriteByte('`')
	case ColumnVal:
		b.sb.WriteByte('?')
		b.addArg(exprTyp.val)
	case Predicate:

		if _, lok := exprTyp.left.(Predicate); lok {
			b.sb.WriteByte('(')
		}

		// 递归左子表达式
		if err := b.buildExpr(exprTyp.left); err != nil {
			return err
		}

		if _, lok := exprTyp.left.(Predicate); lok {
			b.sb.WriteByte(')')
		}

		if exprTyp.left != nil {
			b.sb.WriteByte(' ')
		}
		b.sb.WriteString(string(exprTyp.op))
		if exprTyp.right != nil {
			b.sb.WriteByte(' ')
		}
		if _, rok := exprTyp.right.(Predicate); rok {
			b.sb.WriteByte('(')
		}

		// 递归右子表达式
		if err := b.buildExpr(exprTyp.right); err != nil {
			return err
		}

		if _, rok := exprTyp.right.(Predicate); rok {
			b.sb.WriteByte(')')
		}
	default:
		return errs.UnsupportedExprErr
	}

	return nil
}

func (b *builder) buildCol(col Column) error {
	fd, ok := b.model.Fds[col.name]
	if !ok {
		return errs.InvalidColumnFdErr(col.name)
	}

	b.sb.WriteByte('`')
	b.sb.WriteString(fd.ColName)
	b.sb.WriteByte('`')

	return nil
}

func (b *builder) addArg(val any) {
	if b.args == nil {
		b.args = make([]any, 0, 4)
	}
	b.args = append(b.args, val)
}
