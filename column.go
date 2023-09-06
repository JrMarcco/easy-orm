package orm

type Column struct {
	tbRef   TableRef
	fdName  string
	alias   string
	ufdName string
}

var _ Expression = new(Column)
var _ selectable = new(Column)

func (c Column) expr()       {}
func (c Column) selectable() {}
func (c Column) assign()     {}

// Col 列信息
// 一般作为左子表达式出现。
func Col(fdName string) Column {
	return Column{
		fdName:  fdName,
		ufdName: fdName,
	}
}

func ColWithUpdate(fdName string, ufdName string) Column {
	return Column{
		fdName:  fdName,
		ufdName: ufdName,
	}
}

func (c Column) As(alias string) Column {
	return Column{
		fdName: c.fdName,
		alias:  alias,
	}
}

func (c Column) Eq(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: exprOf(val),
	}
}

func (c Column) Gt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: exprOf(val),
	}
}

func (c Column) Lt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: exprOf(val),
	}
}

// ColumnVal 值信息。
// 一般作为右子表达出现。
type ColumnVal struct {
	val any
}

var _ Expression = new(ColumnVal)

func (v ColumnVal) expr() {}

func valOf(val any) Expression {

	switch valTyp := val.(type) {
	case Predicate:
		return valTyp
	default:
		return ColumnVal{val: val}
	}
}
