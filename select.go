package orm

import (
	"context"
	"strings"
)

type Selector[T any] struct {
	*builder
	conds []condition
	db    *Db
}

func NewSelector[T any](db *Db) *Selector[T] {
	return &Selector[T]{
		builder: newBuilder(),
		db:      db,
	}
}

func (s *Selector[T]) From(tbName string) *Selector[T] {
	s.tbName = tbName
	return s
}

func (s *Selector[T]) Where(predicates ...Predicate) *Selector[T] {
	if s.conds == nil {
		s.conds = make([]condition, 0, 2)
	}

	if len(predicates) > 0 {
		s.conds = append(s.conds, newCond(condTypWhere, predicates))
	}
	return s
}

func (s *Selector[T]) Build() (*Statement, error) {

	var err error
	if s.model, err = s.db.registry.getModel(new(T)); err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT * FROM ")

	if s.tbName == "" {
		s.sb.WriteByte('`')
		s.sb.WriteString(s.model.tbName)
		s.sb.WriteByte('`')
	} else {

		segs := strings.SplitN(s.tbName, ".", 2)

		s.sb.WriteByte('`')
		s.sb.WriteString(segs[0])
		s.sb.WriteByte('`')

		if len(segs) > 1 {
			s.sb.WriteByte('.')
			s.sb.WriteByte('`')
			s.sb.WriteString(segs[1])
			s.sb.WriteByte('`')
		}

	}

	if len(s.conds) > 0 {
		for _, cond := range s.conds {
			s.sb.WriteString(string(cond.typ))

			if err := s.buildExpr(cond.rootExpr); err != nil {
				return nil, err
			}
		}
	}

	s.sb.WriteByte(';')

	return &Statement{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
