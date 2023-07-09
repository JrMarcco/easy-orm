package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/model"
	"reflect"
)

type Inserter[T any] struct {
	*builder
	colFds []string
	rows   []*T
	db     *DB
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		builder: newBuilder(),
		db:      db,
	}
}

func (i *Inserter[T]) ColFd(colFds ...string) *Inserter[T] {
	i.colFds = colFds
	return i
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

	seqFds := i.model.SeqFds

	// 用户指定插入列
	if len(i.colFds) > 0 {
		seqFds = make([]*model.Field, 0, len(i.colFds))
		for _, colFd := range i.colFds {
			fd, ok := i.model.Fds[colFd]
			if !ok {
				return nil, errs.InvalidColumnFdErr(colFd)
			}

			seqFds = append(seqFds, fd)
		}
	}

	i.sb.WriteString("INSERT INTO ")
	i.writeTbName()

	i.sb.WriteByte('(')

	for idx, fd := range seqFds {

		if idx > 0 {
			i.sb.WriteByte(',')
		}

		i.sb.WriteByte('`')
		i.sb.WriteString(fd.ColName)
		i.sb.WriteByte('`')
	}

	i.sb.WriteString(") VALUES ")

	args := make([]any, 0, len(seqFds)*len(i.rows))

	for rowIdx, row := range i.rows {

		if rowIdx > 0 {
			i.sb.WriteByte(',')
		}

		i.sb.WriteByte('(')

		for fdIdx, fd := range seqFds {

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
