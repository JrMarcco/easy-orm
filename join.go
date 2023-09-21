package orm

const (
	JoinTyp      = "JOIN"
	LeftJoinTyp  = "LEFT JOIN"
	RightJoinTyp = "RIGHT JOIN"
)

type TableRef interface {
	table()
}

type Table struct {
	entity any
	alias  string
}

func TableOf(entity any) Table {
	return Table{entity: entity}
}

func TableAs(entity any, alias string) Table {
	return Table{entity: entity, alias: alias}
}

var _ TableRef = new(Table)

func (t Table) table() {}

func (t Table) Col(fdName string) Column {
	return Column{
		tbRef:  t,
		fdName: fdName,
	}
}

func (t Table) Join(target TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTyp,
		left:  t,
		right: target,
	}
}

func (t Table) LeftJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   LeftJoinTyp,
		left:  t,
		right: right,
	}
}

func (t Table) RightJoin(left TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   RightJoinTyp,
		left:  left,
		right: t,
	}
}

type Join struct {
	typ   string
	left  TableRef
	right TableRef
	on    []Predicate
	using []Column
}

var _ TableRef = new(Join)

func (j Join) table() {}

func (j Join) Join(target TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTyp,
		left:  j,
		right: target,
	}
}

func (j Join) LeftJoin(target TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   LeftJoinTyp,
		left:  j,
		right: target,
	}
}

func (j Join) RightJoin(target TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   RightJoinTyp,
		left:  target,
		right: j,
	}
}

type JoinBuilder struct {
	typ   string
	left  TableRef
	right TableRef
}

func (j *JoinBuilder) On(predicates ...Predicate) Join {
	return Join{
		typ:   j.typ,
		left:  j.left,
		right: j.right,
		on:    predicates,
	}
}

func (j *JoinBuilder) Using(cols ...Column) Join {
	return Join{
		typ:   j.typ,
		left:  j.left,
		right: j.right,
		using: cols,
	}
}
