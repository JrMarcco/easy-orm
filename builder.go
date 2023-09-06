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

	registry model.Registry
	dialect  Dialect
	quote    byte
}

func newBuilder(session Session) builder {
	return builder{
		sb:       strings.Builder{},
		registry: session.getCore().registry,
		dialect:  session.getCore().dialect,
		quote:    session.getCore().dialect.quote(),
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
		return errs.ErrInvalidColumnFd(fdName)
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
		b.addArg(exprTyp.val)
		b.dialect.bindArg(b)
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
		return errs.ErrUnsupportedExpr
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
		return errs.ErrUnsupportedSelectable
	}

	return nil
}

func (b *builder) buildCol(col Column) error {
	switch tbRefTyp := col.tbRef.(type) {
	case nil:
		if err := b.writeField(col.fdName); err != nil {
			return err
		}

		if col.alias != "" {
			b.sb.WriteString(" AS ")
			b.writeQuote(col.alias)
		}
	case Table:
		if tbRefTyp.alias != "" {
			b.writeQuote(tbRefTyp.alias)
			b.sb.WriteByte('.')
		}

		m, err := b.registry.Get(tbRefTyp.entity)
		if err != nil {
			return err
		}

		fd, ok := m.Fds[col.fdName]
		if !ok {
			return errs.ErrInvalidColumnFd(col.fdName)
		}
		b.writeQuote(fd.ColName)
	default:
		return errs.ErrInvalidTbRefType(tbRefTyp)
	}

	return nil
}

func (b *builder) buildAggregate(ag Aggregate) error {
	fd, ok := b.model.Fds[ag.fdName]
	if !ok {
		return errs.ErrInvalidColumnFd(ag.fdName)
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

func (b *builder) addArg(vals ...any) {
	if len(vals) == 0 {
		return
	}

	if b.args == nil {
		b.args = make([]any, 0, 4)
	}
	b.args = append(b.args, vals...)
}
