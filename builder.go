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
		if err := b.buildColumn(exprTyp.tableRef, exprTyp.fieldName); err != nil {
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

func (b *builder) buildColumn(tableRef TableRef, fieldName string) error {
	var tableAlias string
	if tableRef != nil {
		tableAlias = tableRef.tableAlias()
	}

	if tableAlias != "" {
		b.writeWithQuote(tableAlias)
		b.sqlBuffer.WriteByte('.')
	}

	columnName, err := b.columnName(tableRef, fieldName)
	if err != nil {
		return err
	}
	b.writeWithQuote(columnName)
	return nil
}

func (b *builder) columnName(tableRef TableRef, fieldName string) (string, error) {
	switch refTyp := tableRef.(type) {
	case nil:
		field, ok := b.model.Fields[fieldName]
		if !ok {
			return "", errs.ErrInvalidField(fieldName)
		}
		return field.ColumnName, nil
	case Table:
		m, err := b.registry.GetModel(refTyp.entity)
		if err != nil {

		}

		field, ok := m.Fields[fieldName]
		if !ok {
			return "", errs.ErrInvalidField(fieldName)
		}
		return field.ColumnName, nil
	case Join:
		columnName, err := b.columnName(refTyp.left, fieldName)
		if err != nil {
			return b.columnName(refTyp.right, fieldName)
		}
		return columnName, nil
	}
	return "", errs.ErrUnsupportedExpr(tableRef)
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

func newBuilder(orm orm) builder {
	dialect := orm.getCore().dialect
	return builder{
		registry:  orm.getCore().registry,
		dialect:   dialect,
		quote:     dialect.quote(),
		sqlBuffer: strings.Builder{},
	}
}
