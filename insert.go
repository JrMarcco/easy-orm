package easyorm

import (
	"context"

	"github.com/JrMarcco/easy-orm/internal/errs"
)

var _ Executor[any] = (*Inserter[any])(nil)

type Inserter[T any] struct {
	builder

	orm  orm
	rows []*T
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	//TODO implement me
	panic("implement me")
}

func (i *Inserter[T]) Insert(rows ...*T) *Inserter[T] {
	i.rows = rows
	return i
}

func (i *Inserter[T]) Build() (*Statement, error) {
	var err error
	if i.model, err = i.orm.getCore().registry.GetModel(new(T)); err != nil {
		return nil, err
	}

	i.sqlBuffer.WriteString("INSERT INTO ")
	i.writeTable()

	if err = i.buildInsertColumns(); err != nil {
		return nil, err
	}

	i.sqlBuffer.WriteByte(';')

	return &Statement{
		SQL:  i.sqlBuffer.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) buildInsertColumns() error {
	if len(i.rows) == 0 {
		return errs.ErrInsertWithoutRows
	}
	fields := i.model.SeqFields

	i.sqlBuffer.WriteString(" (")
	for index, field := range fields {
		if index > 0 {
			i.sqlBuffer.WriteString(", ")
		}

		i.writeWithQuote(field.ColumnName)
	}
	i.sqlBuffer.WriteString(") VALUES ")

	args := make([]any, len(fields)*len(i.rows))
	for rowIndex, row := range i.rows {
		if rowIndex > 0 {
			i.sqlBuffer.WriteString(", ")
		}

		i.sqlBuffer.WriteByte('(')

		resolver := i.orm.getCore().resolverCreator(i.model, row)
		for fieldIndex, field := range fields {
			if fieldIndex > 0 {
				i.sqlBuffer.WriteString(", ")
			}

			val, err := resolver.ReadColumn(field.FiledName)
			if err != nil {
				return err
			}

			args = append(args, val)
			i.dialect.bindArg(&i.builder)
		}
		i.sqlBuffer.WriteByte(')')
	}

	return nil
}

func NewInserter[T any](orm orm) *Inserter[T] {
	return &Inserter[T]{
		builder: newBuilder(orm),
		orm:     orm,
	}
}
