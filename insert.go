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
	args   []any
	db     *DB

	duplicateKey *OnDuplicateKey
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

func (i *Inserter[T]) OnDuplicateKey() *OnDuplicateBuilder[T] {
	return &OnDuplicateBuilder[T]{
		inserter: i,
	}
}

func (i *Inserter[T]) Build() (*Statement, error) {

	if len(i.rows) == 0 {
		return nil, errs.EmptyInsertRowErr
	}

	var err error
	if i.model, err = i.db.registry.Get(new(T)); err != nil {
		return nil, err
	}

	i.sb.WriteString("INSERT INTO ")
	i.writeTbName()

	if err = i.buildInsertCol(); err != nil {
		return nil, err
	}

	if err = i.buildOnDuplicateKey(); err != nil {
		return nil, err
	}

	i.sb.WriteByte(';')

	return &Statement{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) buildInsertCol() error {
	seqFds := i.model.SeqFds

	// 用户指定插入列
	if len(i.colFds) > 0 {
		seqFds = make([]*model.Field, 0, len(i.colFds))
		for _, colFd := range i.colFds {
			fd, ok := i.model.Fds[colFd]
			if !ok {
				return errs.InvalidColumnFdErr(colFd)
			}

			seqFds = append(seqFds, fd)
		}
	}

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

	i.args = make([]any, 0, len(seqFds)*len(i.rows))

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
			i.args = append(i.args, rowVal)
		}
		i.sb.WriteByte(')')
	}

	return nil
}

func (i *Inserter[T]) buildOnDuplicateKey() error {
	if i.duplicateKey != nil {

		i.sb.WriteString(" ON DUPLICATE KEY UPDATE ")

		for idx, assignable := range i.duplicateKey.onConflict {
			if idx > 0 {
				i.sb.WriteByte(',')
			}

			switch typ := assignable.(type) {
			case Assignment:
				fd, ok := i.model.Fds[typ.fdName]
				if !ok {
					return errs.InvalidColumnFdErr(typ.fdName)
				}

				i.sb.WriteByte('`')
				i.sb.WriteString(fd.ColName)
				i.sb.WriteString("`=?")

				i.args = append(i.args, typ.val)
			case Column:
				typ.alias = ""
				if err := i.buildCol(typ); err != nil {
					return err
				}

				i.sb.WriteString("=VALUES(`")

				ufd, ok := i.model.Fds[typ.ufdName]
				if !ok {
					return errs.InvalidColumnFdErr(typ.ufdName)
				}

				i.sb.WriteString(ufd.ColName)
				i.sb.WriteString("`)")
			default:
				return errs.InvalidAssignmentErr
			}
		}
	}

	return nil
}

type OnDuplicateKey struct {
	onConflict []Assignable
}

type OnDuplicateBuilder[T any] struct {
	inserter *Inserter[T]
}

func (o *OnDuplicateBuilder[T]) Update(onConflict ...Assignable) *Inserter[T] {
	o.inserter.duplicateKey = &OnDuplicateKey{
		onConflict: onConflict,
	}
	return o.inserter
}
