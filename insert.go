package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"reflect"
)

type Inserter[T any] struct {
	*builder
	rows []*T
	db   *DB
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		builder: newBuilder(),
		db:      db,
	}
}

func (i *Inserter[T]) Row(rows ...*T) *Inserter[T] {
	i.rows = rows
	return i
}

func (i *Inserter[T]) Build() (*Statement, error) {

	if len(i.rows) == 0 {
		return nil, errs.EmptyInserRowErr
	}

	var err error
	if i.model, err = i.db.registry.Get(new(T)); err != nil {
		return nil, err
	}

	i.sb.WriteString("INSERT INTO ")
	i.writeTbName()

	i.sb.WriteByte('(')

	for idx, fd := range i.model.SeqFds {

		if idx > 0 {
			i.sb.WriteByte(',')
		}

		i.sb.WriteByte('`')
		i.sb.WriteString(fd.ColName)
		i.sb.WriteByte('`')
	}

	i.sb.WriteString(") VALUES ")

	args := make([]any, 0, len(i.model.SeqFds))

	for rowIdx, row := range i.rows {

		if rowIdx > 0 {
			i.sb.WriteByte(',')
		}

		i.sb.WriteByte('(')

		for fdIdx, fd := range i.model.SeqFds {

			if fdIdx > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')

			rowVal := reflect.ValueOf(row).Elem().FieldByName(fd.Name).Interface()
			args = append(args, rowVal)
		}
		i.sb.WriteByte(')')

	}

	i.sb.WriteByte(';')

	return &Statement{
		SQL:  i.sb.String(),
		Args: args,
	}, nil
}
