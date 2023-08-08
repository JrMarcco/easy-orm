package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/model"
	"strings"
)

type builder struct {
	tbName string
	model  *model.Model
	sb     strings.Builder
	args   []any

	dialect Dialect
	quote   byte
}

func newBuilder(dialect Dialect) builder {
	return builder{
		sb:      strings.Builder{},
		dialect: dialect,
		quote:   dialect.quote(),
	}
}

func (b *builder) writeQuote(name string) {
	b.sb.WriteByte(b.quote)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quote)
}

func (b *builder) writeTbName() {
	if b.tbName == "" {
		b.writeQuote(b.model.Tb)
		return
	}

	segs := strings.SplitN(b.tbName, ".", 2)

	b.writeQuote(segs[0])

	if len(segs) > 1 {
		b.sb.WriteByte('.')
		b.writeQuote(segs[1])
	}
}

func (b *builder) writeField(fdName string) error {
	fd, ok := b.model.Fds[fdName]
	if !ok {
		return errs.InvalidColumnFdErr(fdName)
	}
	b.writeQuote(fd.ColName)

	return nil
}

//	buildExpr 构建表达式
//
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

			// 递归右子表达式
			if err := b.buildExpr(exprTyp.right); err != nil {
				return err
			}

			if _, rok := exprTyp.right.(Predicate); rok {
				b.sb.WriteByte(')')
			}
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
	if err := b.writeField(col.fdName); err != nil {
		return err
	}

	if col.alias != "" {
		b.sb.WriteString(" AS ")
		b.writeQuote(col.alias)
	}

	return nil
}

func (b *builder) buildAggregate(ag Aggregate) error {
	fd, ok := b.model.Fds[ag.fdName]
	if !ok {
		return errs.InvalidColumnFdErr(ag.fdName)
	}

	b.sb.WriteString(ag.fnName)

	b.sb.WriteByte('(')
	b.writeQuote(fd.ColName)
	b.sb.WriteByte(')')

	if ag.alias != "" {
		b.sb.WriteString(" AS ")
		b.writeQuote(ag.alias)
	}

	return nil
}

func (b *builder) buildAssign(assign Assignment) error {
	if err := b.writeField(assign.fdName); err != nil {
		return err
	}

	b.sb.WriteString("=?")
	b.addArg(assign.val)

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
