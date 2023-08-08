package orm

import (
	"context"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/model"
	"reflect"
)

type Inserter[T any] struct {
	builder
	colFds     []string
	rows       []*T
	db         *DB
	onConflict *OnConflict
}

func NewInserter[T any](db *DB) *Inserter[T] {
	return &Inserter[T]{
		builder: newBuilder(db.dialect),
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

func (i *Inserter[T]) OnConflicts(conflicts ...string) *OnConflictBuilder[T] {
	b := &OnConflictBuilder[T]{
		inserter: i,
	}

	if len(conflicts) != 0 {
		b.conflicts = conflicts
	}

	return b
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

	if i.onConflict != nil {
		if err = i.dialect.onConflict(&i.builder, i.onConflict); err != nil {
			return nil, err
		}
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

		i.writeQuote(fd.ColName)
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

func (i *Inserter[T]) Exec(ctx context.Context) Result {

	stat, err := i.Build()
	if err != nil {
		return Result{err: err}
	}

	res, err := i.db.sqlDB.ExecContext(ctx, stat.SQL, stat.Args...)
	if err != nil {
		return Result{err: err}
	}

	return Result{
		res: res,
	}
}

type OnConflictBuilder[T any] struct {
	conflicts []string
	inserter  *Inserter[T]
}

func (o *OnConflictBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.inserter.onConflict = &OnConflict{
		conflicts: o.conflicts,
		assigns:   assigns,
	}
	return o.inserter
}
