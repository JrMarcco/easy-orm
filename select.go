package orm

import (
	"context"
	"database/sql"
	"strconv"
)

// selectable 标记接口
// 用来标识可选的查询列
type selectable interface {
	selectable()
}

type Selector[T any] struct {
	builder
	sas    []selectable
	conds  []condition
	limit  int64
	offset int64
	db     *DB
}

var _ Querier[any] = new(Selector[any])

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		builder: newBuilder(db.dialect),
		limit:   0,
		offset:  -1,
		db:      db,
	}
}

func (s *Selector[T]) Select(sas ...selectable) *Selector[T] {
	s.sas = sas
	return s
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

func (s *Selector[T]) GroupBy() *Selector[T] {
	panic("implement me")
}

func (s *Selector[T]) Having() *Selector[T] {
	panic("implement me")
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
	if s.model, err = s.db.registry.Get(new(T)); err != nil {
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
	s.writeTbName()

	if len(s.conds) > 0 {
		for _, cond := range s.conds {
			s.sb.WriteString(string(cond.typ))

			if err := s.buildExpr(cond.rootExpr); err != nil {
				return nil, err
			}
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

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {

	rows, err := s.getRows(ctx)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	res := new(T)

	writer := s.db.creator(s.model, res)
	if err = writer.WriteCols(rows); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	rows, err := s.getRows(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*T, 0, 8)
	for rows.Next() {
		val := new(T)
		writer := s.db.creator(s.model, val)
		if err := writer.WriteCols(rows); err != nil {
			return nil, err
		}

		res = append(res, val)
	}
	return res, nil
}

func (s *Selector[T]) getRows(ctx context.Context) (*sql.Rows, error) {
	stat, err := s.Build()
	if err != nil {
		return nil, err
	}

	sqlDB := s.db.sqlDB

	return sqlDB.QueryContext(ctx, stat.SQL, stat.Args...)
}
