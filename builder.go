package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"strings"
)

type builder struct {
	tbName string
	model  *model
	sb     *strings.Builder
	args   []any
}

// 构建表达式。
// 该过程本是上是一个深度优先遍历二叉树的过程。
func (b *builder) buildExpr(expr Expression) error {
	if expr == nil {
		return nil
	}

	switch exprTyp := expr.(type) {
	case Column:

		fd, ok := b.model.fds[exprTyp.name]
		if !ok {
			return errs.InvalidColumnErr(exprTyp.name)
		}

		b.sb.WriteByte('`')
		b.sb.WriteString(fd.colName)
		b.sb.WriteByte('`')
	case Value:
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

func (b *builder) addArg(val any) {
	if b.args == nil {
		b.args = make([]any, 0, 4)
	}
	b.args = append(b.args, val)
}
