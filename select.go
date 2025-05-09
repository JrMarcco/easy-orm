package easyorm

import (
	"context"
)

// selectable marker interface, used to identify optional query columns( e.g., columns, aggregate functions ).
type selectable interface {
	selectable()
}

var _ Querier[any] = (*Selector[any])(nil)

type Selector[T any] struct {
	builder
	session

	limit  int64
	offset int64

	where []Predicate
}

func (s *Selector[T]) FindOne(ctx context.Context) (T, error) {
	// TODO implement me
	panic("implement me")
}

func (s *Selector[T]) FindMulti(ctx context.Context) ([]T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) Where(pds ...Predicate) *Selector[T] {
	s.where = append(s.where, pds...)
	return s
}

func (s *Selector[T]) Limit(limit int64) *Selector[T] {
	s.limit = limit
	return s
}

func (s *Selector[T]) Offset(offset int64) *Selector[T] {
	s.offset = offset
	return s
}

func (s *Selector[T]) Build() (*Statement, error) {
	var err error
	if s.model, err = s.getCore().registry.GetModel(new(T)); err != nil {
		return nil, err
	}

	s.sqlBuffer.WriteString("SELECT * FROM ")

	s.writeTable()

	lenOfWhere := len(s.where)
	if lenOfWhere > 0 {
		s.sqlBuffer.WriteString(" WHERE ")

		pd := s.where[0]
		for i := 1; i < lenOfWhere; i++ {
			pd = pd.And(s.where[i])
		}
		if err := s.buildExpr(pd); err != nil {
			return nil, err
		}
	}

	s.sqlBuffer.WriteByte(';')
	return &Statement{
		Sql:  s.sqlBuffer.String(),
		Args: s.args,
	}, nil
}

func NewSelector[T any](session session) *Selector[T] {
	return &Selector[T]{
		builder: newBuilder(session),
		session: session,
	}
}
