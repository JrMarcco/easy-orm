package orm

import (
	"context"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"strconv"
)

// selectable 标记接口
// 用来标识可选的查询列（列、聚合函数等）
type selectable interface {
	selectable()
}

type Selector[T any] struct {
	builder

	tbRef       TableRef
	sas         []selectable
	whereConds  []condition
	havingConds []condition
	groupByCols []Column
	limit       int64
	offset      int64

	session Session
}

var _ Querier[any] = new(Selector[any])

func NewSelector[T any](session Session) *Selector[T] {
	return &Selector[T]{
		builder: newBuilder(session),
		limit:   0,
		offset:  -1,
		session: session,
	}
}

func (s *Selector[T]) Select(sas ...selectable) *Selector[T] {
	s.sas = sas
	return s
}

func (s *Selector[T]) From(tbRef TableRef) *Selector[T] {
	s.tbRef = tbRef
	return s
}

func (s *Selector[T]) Where(predicates ...Predicate) *Selector[T] {
	if s.whereConds == nil {
		s.whereConds = make([]condition, 0, 2)
	}

	if len(predicates) > 0 {
		s.whereConds = append(s.whereConds, newCond(condTypWhere, predicates))
	}
	return s
}

func (s *Selector[T]) GroupBy(groupByCols ...Column) *Selector[T] {
	s.groupByCols = groupByCols
	return s
}

func (s *Selector[T]) Having(predicates ...Predicate) *Selector[T] {
	if s.havingConds == nil {
		s.havingConds = make([]condition, 0, 2)
	}

	if len(predicates) > 0 {
		s.havingConds = append(s.havingConds, newCond(condTypHaving, predicates))
	}
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

func (s *Selector[T]) AsSubQuery(alias string) SubQ {
	tbRef := s.tbRef
	if tbRef == nil {
		tbRef = TableOf(new(T))
	}

	return SubQ{
		builder: s,
		sas:     s.sas,
		alias:   alias,
		tbRef:   tbRef,
	}
}

func (s *Selector[T]) Build() (*Statement, error) {

	var err error
	if s.model, err = s.session.getCore().registry.Get(new(T)); err != nil {
		return nil, err
	}

	s.sb.WriteString("SELECT ")

	if len(s.sas) > 0 {
		for i, sa := range s.sas {
			if i > 0 {
				s.sb.WriteByte(',')
			}

			if err = s.buildSelectable(sa); err != nil {
				return nil, err
			}
		}
	} else {
		s.sb.WriteByte('*')
	}

	s.sb.WriteString(" FROM ")
	if err = s.buildTable(s.tbRef); err != nil {
		return nil, err
	}

	if len(s.whereConds) != 0 {
		if err = s.buildConds(s.whereConds); err != nil {
			return nil, err
		}
	}

	if len(s.groupByCols) != 0 {
		s.sb.WriteString(" GROUP BY ")
		for idx, col := range s.groupByCols {

			if idx > 0 {
				s.sb.WriteByte(',')
			}

			if err = s.writeField(col.fdName); err != nil {
				return nil, err
			}
		}
	}

	if len(s.havingConds) != 0 {
		// 校验是否有 group by
		if len(s.groupByCols) == 0 {
			return nil, errs.ErrHavingWithoutGroupBy
		}

		if err = s.buildConds(s.havingConds); err != nil {
			return nil, err
		}
	}

	if s.limit != 0 {
		s.sb.WriteString(" LIMIT ")
		s.sb.WriteString(strconv.FormatInt(s.limit, 10))
	}

	if s.offset != -1 {
		s.sb.WriteString(" OFFSET ")
		s.sb.WriteString(strconv.FormatInt(s.offset, 10))
	}

	s.sb.WriteByte(';')

	return &Statement{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildTable(tbRef TableRef) error {
	switch tbRefTyp := tbRef.(type) {
	case nil:
		s.writeTbName()
	case Table:
		md, err := s.session.getCore().registry.Get(tbRefTyp.entity)
		if err != nil {
			return err
		}
		s.writeQuote(md.Tb)

		if tbRefTyp.alias != "" {
			s.sb.WriteString(" AS ")
			s.writeQuote(tbRefTyp.alias)
		}
	case Join:
		s.sb.WriteByte('(')

		if err := s.buildTable(tbRefTyp.left); err != nil {
			return err
		}

		s.sb.WriteByte(' ')
		s.sb.WriteString(tbRefTyp.typ)
		s.sb.WriteByte(' ')

		if err := s.buildTable(tbRefTyp.right); err != nil {
			return err
		}

		if len(tbRefTyp.using) > 0 {
			s.sb.WriteString(" USING (")

			for i, col := range tbRefTyp.using {
				if i > 0 {
					s.sb.WriteByte(',')
				}

				if err := s.writeField(col.fdName); err != nil {
					return nil
				}
			}

			s.sb.WriteByte(')')
		}

		if len(tbRefTyp.on) > 0 {
			s.sb.WriteString(" ON ")
			for i, pd := range tbRefTyp.on {
				if i > 0 {
					s.sb.WriteString(" AND ")
				}
				if err := s.buildExpr(pd); err != nil {
					return err
				}
			}
		}
		s.sb.WriteByte(')')
	case SubQ:
		return s.buildSubQuery(tbRefTyp)
	default:
		return errs.ErrInvalidTbRefType(tbRefTyp)
	}

	return nil
}

func (s *Selector[T]) buildConds(conds []condition) error {
	for _, cond := range conds {
		s.sb.WriteString(string(cond.typ))

		if err := s.buildExpr(cond.rootExpr); err != nil {
			return err
		}
	}

	return nil
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	if s.limit != 1 {
		s.limit = 1
	}

	return get[T](ctx, s.session, &StatContext{
		Typ:     ScTypSelect,
		Builder: s,
	})
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	return getMulti[T](ctx, s.session, &StatContext{
		Typ:     ScTypSelect,
		Builder: s,
	})
}
