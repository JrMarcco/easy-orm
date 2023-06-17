package orm

type op string

const (
	opEq  = "="
	opGt  = ">"
	opLt  = "<"
	opAnd = "AND"
	opOr  = "OR"
	opNot = "NOT"
)

// Predicate 条件
// 表达式的一种，由左右两个子表达式以及中间的操作符组成。
// 其中部分操作符左子表达式可为空，例如 NOT。
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

func Not(right Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: right,
	}
}
