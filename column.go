package orm

type Column struct {
	name string
}

func (c Column) expr() {}

// Col 列信息
// 一般作为左子表达式出现。
func Col(name string) Column {
	return Column{name: name}
}

func (c Column) Eq(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: valOf(val),
	}
}

func (c Column) Gt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: valOf(val),
	}
}

func (c Column) Lt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: valOf(val),
	}
}

// Value 值信息。
// 一般作为右子表达出现。
type Value struct {
	val any
}

func (v Value) expr() {}

func valOf(val any) Value {
	return Value{val: val}
}
