package easyorm

var _ selectable = (*Column)(nil)
var _ Expression = (*Column)(nil)

type Column struct {
	name string
}

func (c Column) selectable() {}
func (c Column) expr()       {}

func (c Column) Eq(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: ColumnValue{value: val},
	}
}

func (c Column) Ne(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opNe,
		right: ColumnValue{value: val},
	}
}

func (c Column) Gt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: ColumnValue{value: val},
	}
}

func (c Column) Ge(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGe,
		right: ColumnValue{value: val},
	}
}
func (c Column) Lt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: ColumnValue{value: val},
	}
}

func (c Column) Le(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLe,
		right: ColumnValue{value: val},
	}
}

func Col(name string) Column {
	return Column{
		name: name,
	}
}

var _ Expression = (*ColumnValue)(nil)

type ColumnValue struct {
	value any
}

func (c ColumnValue) expr() {}

func valOf(val any) Expression {
	switch valType := val.(type) {
	case Predicate:
		return valType
	default:
		return ColumnValue{
			value: val,
		}
	}
}
