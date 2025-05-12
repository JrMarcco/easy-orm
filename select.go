package easyorm

import (
	"context"
	"strconv"

	"github.com/JrMarcco/easy-orm/internal/errs"
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
	having      []Condition
	groupBy     []Column
	orderBy     []OrderBy
}

func (s *Selector[T]) FindOne(ctx context.Context) (*T, error) {
	if s.limit != 1 {
		s.limit = 1
	}

	return findOne[T](ctx, &OrmContext{
		Typ:     ScTypSELECT,
		Builder: s,
	}, s.orm)
}

func (s *Selector[T]) FindMulti(ctx context.Context) ([]*T, error) {
	return findMulti[T](ctx, &OrmContext{
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

func (s *Selector[T]) GroupBy(cols ...Column) *Selector[T] {
	s.groupBy = cols
	return s
}

func (s *Selector[T]) Having(pds ...Predicate) *Selector[T] {
	if len(pds) == 0 {
		return s
	}

	if s.having == nil {
		s.having = make([]Condition, 0, 4)
	}
	s.having = append(s.having, NewCondition(condTypHaving, pds))
	return s
}

func (s *Selector[T]) OrderBy(orderBys ...OrderBy) *Selector[T] {
	s.orderBy = orderBys
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
		if err = s.buildConditions(s.where); err != nil {
			return nil, err
		}
	}

	if len(s.groupBy) > 0 {
		s.sqlBuffer.WriteString(" GROUP BY ")
		for index, col := range s.groupBy {
			field, ok := s.model.Fields[col.fieldName]
			if !ok {
				return nil, errs.ErrInvalidField(col.fieldName)
			}

			if index > 0 {
				s.sqlBuffer.WriteString(", ")
			}
			s.writeWithQuote(field.ColumnName)
		}
	}

	if len(s.having) > 0 {
		if len(s.groupBy) == 0 {
			return nil, errs.ErrHavingWithoutGroupBy
		}
		if err = s.buildConditions(s.having); err != nil {
			return nil, err
		}
	}

	if len(s.orderBy) > 0 {
		if err = s.buildOrderBy(); err != nil {
			return nil, err
		}
	}

	if s.limit > 0 {
		s.sqlBuffer.WriteString(" LIMIT ")
		s.sqlBuffer.WriteString(strconv.FormatInt(s.limit, 10))
	}

	if s.offset > 0 {
		s.sqlBuffer.WriteString(" OFFSET ")
		s.sqlBuffer.WriteString(strconv.FormatInt(s.offset, 10))
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

func (s *Selector[T]) buildOrderBy() error {
	if len(s.orderBy) == 0 {
		return nil
	}

	s.sqlBuffer.WriteString(" ORDER BY ")
	for index, ob := range s.orderBy {
		if index > 0 {
			s.sqlBuffer.WriteString(", ")
		}

		field, ok := s.model.Fields[ob.fieldName]
		if !ok {
			return errs.ErrInvalidField(ob.fieldName)
		}
		s.writeWithQuote(field.ColumnName)
		s.sqlBuffer.WriteString(" ")
		s.sqlBuffer.WriteString(ob.typ.String())
	}
	return nil
}

type orderTyp string

const (
	orderAsc  orderTyp = "ASC"
	orderDesc orderTyp = "DESC"
)

func (o orderTyp) String() string {
	return string(o)
}

type OrderBy struct {
	fieldName string
	typ       orderTyp
}

func Asc(fieldName string) OrderBy {
	return OrderBy{
		fieldName: fieldName,
		typ:       orderAsc,
	}
}

func Desc(fieldName string) OrderBy {
	return OrderBy{
		fieldName: fieldName,
		typ:       orderDesc,
	}
}

func (s *Selector[T]) buildConditions(conditions []Condition) error {
	for _, c := range conditions {
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
