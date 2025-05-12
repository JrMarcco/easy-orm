package easyorm

import (
	"context"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/model"
)

var _ Executor[any] = (*Inserter[any])(nil)

type Inserter[T any] struct {
	builder

	orm    orm
	rows   []*T
	fields []string

	conflict *Conflict
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	return exec(ctx, &StatementContext{
		Typ:     ScTypINSERT,
		Builder: i,
	}, i.orm)
}

func (i *Inserter[T]) Rows(rows ...*T) *Inserter[T] {
	i.rows = rows
	return i
}

func (i *Inserter[T]) Fields(fields ...string) *Inserter[T] {
	i.fields = fields
	return i
}

// OnConflict upsert support.
// conflicts only supported on postgres, conflicts are field in entity, not columns in db table.
func (i *Inserter[T]) OnConflict(conflicts ...string) *OnConflictBuilder[T] {
	ocb := &OnConflictBuilder[T]{
		inserter: i,
	}

	if len(conflicts) > 0 {
		ocb.conflicts = conflicts
	}
	return ocb
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

	if i.conflict != nil {
		if err = i.dialect.onConflict(&i.builder, i.conflict); err != nil {
			return nil, err
		}
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

	if len(i.fields) > 0 {
		fields = make([]*model.Field, 0, len(i.fields))
		for _, f := range i.fields {
			field, ok := i.model.Fields[f]
			if !ok {
				return errs.ErrInvalidField(f)
			}
			fields = append(fields, field)
		}
	}

	i.sqlBuffer.WriteString(" (")
	for index, field := range fields {
		if index > 0 {
			i.sqlBuffer.WriteString(", ")
		}

		i.writeWithQuote(field.ColumnName)
	}
	i.sqlBuffer.WriteString(") VALUES ")

	i.args = make([]any, 0, len(fields)*len(i.rows))
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
			i.args = append(i.args, val)
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

type OnConflictBuilder[T any] struct {
	inserter  *Inserter[T]
	conflicts []string
}

func (o *OnConflictBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.inserter.conflict = &Conflict{
		conflicts: o.conflicts,
		assigns:   assigns,
	}
	return o.inserter
}
