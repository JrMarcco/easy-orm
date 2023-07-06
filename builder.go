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
	case ColumnVal:
		b.sb.WriteByte('?')
		b.addArg(exprTyp.val)
	case RawExpr:
		b.sb.WriteString(exprTyp.raw)
		b.addArg(exprTyp.args...)
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

		if exprTyp.op != "" {
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

func (b *builder) buildSelectable(sa selectable) error {
	switch saType := sa.(type) {
	case Column:
		if err := b.buildCol(saType); err != nil {
			return err
		}
	case Aggregate:
		if err := b.buildAggregate(saType); err != nil {
			return err
		}
	case RawExpr:
		b.sb.WriteString(saType.raw)
		b.addArg(saType.args...)
	default:
		return errs.UnsupportedSelectableErr
	}

	return nil
}

func (b *builder) buildCol(col Column) error {
	fd, ok := b.model.Fds[col.fdName]
	if !ok {
		return errs.InvalidColumnFdErr(col.fdName)
	}

	b.sb.WriteByte('`')
	b.sb.WriteString(fd.ColName)
	b.sb.WriteByte('`')

	return nil
}

func (b *builder) buildAggregate(ag Aggregate) error {
	fd, ok := b.model.Fds[ag.fdName]
	if !ok {
		return errs.InvalidColumnFdErr(ag.fdName)
	}

	b.sb.WriteString(ag.fnName)
	b.sb.WriteString("(`")
	b.sb.WriteString(fd.ColName)
	b.sb.WriteString("`)")

	return nil
}

func (b *builder) addArg(vals ...any) {
	if len(vals) == 0 {
		return
	}

	if b.args == nil {
		b.args = make([]any, 0, 4)
	}
	b.args = append(b.args, vals...)
}
