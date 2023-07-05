package orm

type Column struct {
	fdName string
}

var _ Expression = new(Column)
var _ selectable = new(Column)

func (c Column) expr() {}

func (c Column) selectable() {}

// Col 列信息
// 一般作为左子表达式出现。
func Col(name string) Column {
	return Column{fdName: name}
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

// ColumnVal 值信息。
// 一般作为右子表达出现。
type ColumnVal struct {
	val any
}

var _ Expression = new(ColumnVal)

func (v ColumnVal) expr() {}

func valOf(val any) ColumnVal {
	return ColumnVal{val: val}
}
