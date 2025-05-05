package easyorm

// Expression marker interface, representing an expression.
// All the elements after "WHERE" are expression.
type Expression interface {
	expr()
}

func exprOf(val any) Expression {
	switch valType := val.(type) {
	case Expression:
		return valType
	case Predicate:
		return valType
	default:
		return ColumnValue{
			value: val,
		}
	}
}
