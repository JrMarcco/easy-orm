package easyorm

var _ selectable = (*Column)(nil)
var _ Expression = (*Column)(nil)

type Column struct {
	fieldName string
	alias     string
}

func (c Column) selectable() {}
func (c Column) expr()       {}

func (c Column) As(alias string) Column {
	return Column{
		fieldName: c.fieldName,
		alias:     alias,
	}
}

func (c Column) Eq(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: valueOf(val),
	}
}

func (c Column) Ne(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opNe,
		right: valueOf(val),
	}
}

func (c Column) Gt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGt,
		right: valueOf(val),
	}
}

func (c Column) Ge(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opGe,
		right: valueOf(val),
	}
}
func (c Column) Lt(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLt,
		right: valueOf(val),
	}
}

func (c Column) Le(val any) Predicate {
	return Predicate{
		left:  c,
		op:    opLe,
		right: valueOf(val),
	}
}

// Col create a column expression.
//
// fieldName is the field name of the model.
func Col(fieldName string) Column {
	return Column{
		fieldName: fieldName,
	}
}

var _ Expression = (*columnValue)(nil)

type columnValue struct {
	value any
}

func (c columnValue) expr() {}

func valueOf(v any) Expression {
	switch typ := v.(type) {
	case Expression:
		return typ
	default:
		return columnValue{value: v}
	}
}
