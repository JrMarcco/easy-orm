package easyorm

import (
	"context"
)

var _ Executor[any] = (*Deleter[any])(nil)

type Deleter[T any] struct {
	builder

	orm   orm
	where []Condition
}

func (d *Deleter[T]) Exec(ctx context.Context) Result {
	return exec(ctx, &StatementContext{
		Typ:     ScTypDELETE,
		Builder: d,
	}, d.orm)
}

func (d *Deleter[T]) Where(pds ...Predicate) *Deleter[T] {
	if len(pds) == 0 {
		return d
	}

	if d.where == nil {
		d.where = make([]Condition, 0, 1)
	}

	d.where = append(d.where, NewCondition(condTypWhere, pds))
	return d
}

func (d *Deleter[T]) Build() (*Statement, error) {
	var err error

	if d.model, err = d.orm.getCore().registry.GetModel(new(T)); err != nil {
		return nil, err
	}

	d.sqlBuffer.WriteString("DELETE FROM ")
	d.writeTable()

	if d.where != nil {
		if err = d.buildCondition(); err != nil {
			return nil, err
		}
	}

	d.sqlBuffer.WriteByte(';')
	return &Statement{
		SQL:  d.sqlBuffer.String(),
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) buildCondition() error {
	for _, c := range d.where {
		d.sqlBuffer.WriteString(c.typ.String())
		if err := d.buildExpr(c.expr); err != nil {
			return err
		}
	}
	return nil
}

func NewDeleter[T any](orm orm) *Deleter[T] {
	return &Deleter[T]{
		builder: newBuilder(orm),
		orm:     orm,
	}
}
