package easyorm

import "fmt"

type condTyp string

func (c condTyp) String() string {
	return fmt.Sprintf(" %s ", string(c))
}

const (
	condTypWhere  condTyp = "WHERE"
	condTypHaving condTyp = "HAVING"
)

type Condition struct {
	typ  condTyp
	expr Expression
}

func NewCondition(typ condTyp, pds []Predicate) Condition {
	expr := pds[0]

	for _, pd := range pds[1:] {
		expr = expr.And(pd)
	}
	return Condition{
		typ:  typ,
		expr: expr,
	}
}

type op string

func (o op) String() string {
	return string(o)
}

const (
	opEq  op = "="
	opGt  op = ">"
	opLt  op = "<"
	opGe  op = ">="
	opLe  op = "<="
	opNe  op = "!="
	opAnd op = "AND"
	opOr  op = "OR"
	opNot op = "NOT"
	opIn  op = "IN"
)

var _ Expression = (*Predicate)(nil)

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

func (p Predicate) expr() {}

func (p Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opAnd,
		right: right,
	}
}

func (p Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  p,
		op:    opOr,
		right: right,
	}
}

func (p Predicate) Not() Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}
