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

	orm orm

	limit  int64
	offset int64

	selectables []selectable
	where       []Condition
}

func (s *Selector[T]) FindOne(ctx context.Context) (*T, error) {
	if s.limit != 1 {
		s.limit = 1
	}

	return findOne[T](ctx, &StatementContext{
		Typ:     ScTypSELECT,
		Builder: s,
	}, s.orm)
}

func (s *Selector[T]) FindMulti(ctx context.Context) ([]*T, error) {
	return findMulti[T](ctx, &StatementContext{
		Typ:     ScTypSELECT,
		Builder: s,
	}, s.orm)
}

func (s *Selector[T]) Select(selectables ...selectable) *Selector[T] {
	s.selectables = selectables
	return s
}

func (s *Selector[T]) Where(pds ...Predicate) *Selector[T] {
	if len(pds) == 0 {
		return s
	}

	if s.where == nil {
		s.where = make([]Condition, 0, 4)
	}

	s.where = append(s.where, NewCondition(condTypWhere, pds))
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
	if s.model, err = s.orm.getCore().registry.GetModel(new(T)); err != nil {
		return nil, err
	}

	s.sqlBuffer.WriteString("SELECT ")
	if err = s.buildSelectedColumns(); err != nil {
		return nil, err
	}
	s.sqlBuffer.WriteString(" FROM ")

	s.writeTable()

	if len(s.where) > 0 {
		if err = s.buildConditions(); err != nil {
			return nil, err
		}
	}

	s.sqlBuffer.WriteByte(';')
	return &Statement{
		SQL:  s.sqlBuffer.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildSelectedColumns() error {
	if len(s.selectables) == 0 {
		s.sqlBuffer.WriteByte('*')
		return nil
	}

	for i, sa := range s.selectables {
		if i > 0 {
			s.sqlBuffer.WriteString(", ")
		}

		if err := s.buildSelectable(sa); err != nil {
			return err
		}
	}

	return nil
}

func (s *Selector[T]) buildConditions() error {
	for _, c := range s.where {
		s.sqlBuffer.WriteString(c.typ.String())

		if err := s.buildExpr(c.expr); err != nil {
			return err
		}
	}
	return nil
}

var _ StatementBuilder = (*Selector[any])(nil)

func NewSelector[T any](orm orm) *Selector[T] {
	return &Selector[T]{
		builder: newBuilder(orm),
		orm:     orm,
		limit:   0,
		offset:  -1,
	}
}
