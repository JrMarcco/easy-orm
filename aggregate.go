package orm

type Aggregate struct {
	fnName string
	fdName string
	alias  string
}

var _ selectable = new(Aggregate)

func (a Aggregate) selectable() {}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		fnName: a.fnName,
		fdName: a.fdName,
		alias:  alias,
	}
}

func Avg(fdName string) Aggregate {
	return Aggregate{
		fnName: "AVG",
		fdName: fdName,
	}
}

func Sum(fdName string) Aggregate {
	return Aggregate{
		fnName: "SUM",
		fdName: fdName,
	}
}

func Count(fdName string) Aggregate {
	return Aggregate{
		fnName: "COUNT",
		fdName: fdName,
	}
}

func Max(fdName string) Aggregate {
	return Aggregate{
		fnName: "MAX",
		fdName: fdName,
	}
}

func Min(fdName string) Aggregate {
	return Aggregate{
		fnName: "MIN",
		fdName: fdName,
	}
}
