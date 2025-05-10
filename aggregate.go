package easyorm

var _ selectable = (*Aggregate)(nil)

type Aggregate struct {
	funcName  string
	fieldName string
	alias     string
}

func (a Aggregate) selectable() {}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		funcName:  a.funcName,
		fieldName: a.fieldName,
		alias:     alias,
	}
}

func Count(fieldName string) Aggregate {
	return Aggregate{
		funcName:  "COUNT",
		fieldName: fieldName,
	}
}

func Sum(fieldName string) Aggregate {
	return Aggregate{
		funcName:  "SUM",
		fieldName: fieldName,
	}
}

func Max(fieldName string) Aggregate {
	return Aggregate{
		funcName:  "MAX",
		fieldName: fieldName,
	}
}

func Min(fieldName string) Aggregate {
	return Aggregate{
		funcName:  "MIN",
		fieldName: fieldName,
	}
}

func Avg(fieldName string) Aggregate {
	return Aggregate{
		funcName:  "AVG",
		fieldName: fieldName,
	}
}
