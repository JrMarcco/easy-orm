package orm

import (
	"context"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/model"
)

type Inserter[T any] struct {
	builder
	colFds     []string
	rows       []*T
	onConflict *OnConflict

	session Session
}

var _ Executor = &Inserter[any]{}

func NewInserter[T any](session Session) *Inserter[T] {
	return &Inserter[T]{
		builder: newBuilder(session),
		session: session,
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
		return nil, errs.ErrEmptyInsertRow
	}

	var err error
	if i.model, err = i.session.getCore().registry.Get(new(T)); err != nil {
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
				return errs.ErrInvalidColumnFd(colFd)
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

		valCreator := i.session.getCore().creator(i.model, row)
		for fdIdx, fd := range seqFds {
			rowVal, err := valCreator.ReadCol(fd.Name)
			if err != nil {
				return err
			}
			i.args = append(i.args, rowVal)

			if fdIdx > 0 {
				i.sb.WriteByte(',')
			}
			i.dialect.bindArg(&i.builder)
		}
		i.sb.WriteByte(')')
	}

	return nil
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	return exec(ctx, i.session, &StatContext{
		Typ:     ScTypInsert,
		Builder: i,
	})
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
