package easyorm

import (
	"strings"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/model"
)

type builder struct {
	model    *model.Model
	registry model.Registry

	dialect Dialect
	quote   byte

	sqlBuffer strings.Builder
	args      []any
}

func (b *builder) writeWithQuote(name string) {
	b.sqlBuffer.WriteByte(b.quote)
	b.sqlBuffer.WriteString(name)
	b.sqlBuffer.WriteByte(b.quote)
}

func (b *builder) writeTable() {
	segments := strings.Split(b.model.TableName, ".")
	b.writeWithQuote(segments[0])

	if len(segments) > 1 {
		b.writeWithQuote(segments[1])
	}
}

func (b *builder) writeField(name string) error {
	if field, ok := b.model.Fields[name]; ok {
		b.writeWithQuote(field.ColumnName)
		return nil
	}
	return errs.ErrInvalidField(name)
}

func (b *builder) buildExpr(expr Expression) error {
	if expr == nil {
		return nil
	}

	switch exprTyp := expr.(type) {
	case Predicate:
		leftIsPredicate := false
		if _, ok := exprTyp.left.(Predicate); ok {
			leftIsPredicate = true
		}

		if leftIsPredicate {
			b.sqlBuffer.WriteByte('(')
		}
		if err := b.buildExpr(exprTyp.left); err != nil {
			return err
		}
		if leftIsPredicate {
			b.sqlBuffer.WriteByte(')')
		}

		if exprTyp.op != "" {
			if exprTyp.left != nil {
				b.sqlBuffer.WriteByte(' ')
			}
			b.sqlBuffer.WriteString(string(exprTyp.op))
			if exprTyp.right != nil {
				b.sqlBuffer.WriteByte(' ')
			}

			rightIsPredicate := false
			if _, ok := exprTyp.right.(Predicate); ok {
				rightIsPredicate = true
			}
			if rightIsPredicate {
				b.sqlBuffer.WriteByte('(')
			}
			if err := b.buildExpr(exprTyp.right); err != nil {
				return err
			}
			if rightIsPredicate {
				b.sqlBuffer.WriteByte(')')
			}
		}
	case Column:
		if err := b.buildColumn(exprTyp); err != nil {
			return err
		}
	case columnValue:
		b.addArgs(exprTyp.value)
		b.dialect.bindArg(b)
	case RawExpression:
		b.sqlBuffer.WriteString(exprTyp.raw)
		b.addArgs(exprTyp.args...)
	default:
		return errs.ErrUnsupportedExpr(exprTyp)
	}
	return nil
}

func (b *builder) buildSelectable(sa selectable) error {
	switch saTyp := sa.(type) {
	case Column:
		return b.buildColumn(saTyp)
	case Aggregate:
		return b.buildAggregate(saTyp)
	}
	return nil
}

func (b *builder) buildColumn(column Column) error {
	if err := b.writeField(column.fieldName); err != nil {
		return err
	}

	if column.alias != "" {
		b.sqlBuffer.WriteString(" AS ")
		b.writeWithQuote(column.alias)
	}

	return nil
}

func (b *builder) buildAggregate(aggregate Aggregate) error {
	field, ok := b.model.Fields[aggregate.fieldName]
	if !ok {
		return errs.ErrInvalidField(aggregate.fieldName)
	}

	b.sqlBuffer.WriteString(aggregate.funcName)
	b.sqlBuffer.WriteByte('(')
	b.writeWithQuote(field.ColumnName)
	b.sqlBuffer.WriteByte(')')

	if aggregate.alias != "" {
		b.sqlBuffer.WriteString(" AS ")
		b.writeWithQuote(aggregate.alias)
	}

	return nil
}

func (b *builder) addArgs(val ...any) {
	if len(val) == 0 {
		return
	}
	if b.args == nil {
		b.args = make([]any, 0, 8)
	}

	b.args = append(b.args, val...)
}

func newBuilder(session orm) builder {
	dialect := session.getCore().dialect
	return builder{
		registry:  session.getCore().registry,
		dialect:   dialect,
		quote:     dialect.quote(),
		sqlBuffer: strings.Builder{},
	}
}
