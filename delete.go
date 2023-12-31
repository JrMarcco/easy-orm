package orm

import (
	"context"
	"strings"
)

type Deleter[T any] struct {
	builder
	conds []condition

	session Session
}

var _ Executor = &Deleter[any]{}

func NewDeleter[T any](session Session) *Deleter[T] {
	return &Deleter[T]{
		builder: newBuilder(session),
		session: session,
	}
}

func (d *Deleter[T]) From(tbName string) *Deleter[T] {
	d.tbName = tbName
	return d
}

func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	if d.conds == nil {
		d.conds = make([]condition, 0, 2)
	}

	if len(predicates) > 0 {
		d.conds = append(d.conds, newCond(condTypWhere, predicates))
	}
	return d
}

func (d *Deleter[T]) Build() (*Statement, error) {

	var err error
	if d.model, err = d.session.getCore().registry.Get(new(T)); err != nil {
		return nil, err
	}

	d.sb = strings.Builder{}
	d.sb.WriteString("DELETE FROM ")
	d.writeTbName()

	if len(d.conds) > 0 {
		for _, cond := range d.conds {
			d.sb.WriteString(string(cond.typ))

			if err := d.buildExpr(cond.rootExpr); err != nil {
				return nil, err
			}
		}
	}

	d.sb.WriteByte(';')

	return &Statement{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) Exec(ctx context.Context) Result {
	return exec(ctx, d.session, &StatContext{
		Typ:     ScTypDelete,
		Builder: d,
	})
}
