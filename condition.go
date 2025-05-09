package easyorm

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
