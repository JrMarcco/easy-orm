package easyorm

type joinTyp string

func (j joinTyp) String() string {
	return string(j)
}

const (
	JoinTypeInner joinTyp = "INNER JOIN"
	JoinTypeLeft  joinTyp = "LEFT JOIN"
	JoinTypeRight joinTyp = "RIGHT JOIN"
)

type TableRef interface {
	tableAlias() string
}

var _ TableRef = (*Table)(nil)

type Table struct {
	entity any
	alias  string
}

func (t Table) tableAlias() string {
	return t.alias
}

func (t Table) As(alias string) Table {
	return Table{
		entity: t.entity,
		alias:  alias,
	}
}

func (t Table) Col(fieldName string) Column {
	return Column{
		tableRef:  t,
		fieldName: fieldName,
	}
}

func (t Table) InnerJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeInner,
		left:  t,
		right: right,
	}
}

func (t Table) LeftJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeLeft,
		left:  t,
		right: right,
	}
}

func (t Table) RightJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeRight,
		left:  t,
		right: right,
	}
}

func TableOf(entity any) Table {
	return Table{
		entity: entity,
	}
}

func TableAs(entity any, alias string) Table {
	return Table{
		entity: entity,
		alias:  alias,
	}
}

var _ TableRef = (*Join)(nil)

type Join struct {
	typ   joinTyp
	left  TableRef
	right TableRef
	on    []Predicate
	using []Column
}

func (j Join) tableAlias() string {
	return ""
}

func (j Join) InnerJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeInner,
		left:  j,
		right: right,
	}
}

func (j Join) LeftJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeLeft,
		left:  j,
		right: right,
	}
}

func (j Join) RightJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeRight,
		left:  j,
		right: right,
	}
}

type JoinBuilder struct {
	typ   joinTyp
	left  TableRef
	right TableRef
}

func (b *JoinBuilder) On(pds ...Predicate) Join {
	return Join{
		typ:   b.typ,
		left:  b.left,
		right: b.right,
		on:    pds,
	}
}

func (b *JoinBuilder) Using(cols ...Column) Join {
	return Join{
		typ:   b.typ,
		left:  b.left,
		right: b.right,
		using: cols,
	}
}

var _ selectable = (*SubQuery)(nil)
var _ Expression = (*SubQuery)(nil)
var _ TableRef = (*SubQuery)(nil)

type SubQuery struct {
	builder StatementBuilder

	tableRef    TableRef
	selectables []selectable

	alias string
}

func (s SubQuery) selectable() {}
func (s SubQuery) expr()       {}

func (s SubQuery) tableAlias() string {
	return s.alias
}

func (s SubQuery) As(alias string) SubQuery {
	return SubQuery{
		builder:     s.builder,
		tableRef:    s.tableRef,
		selectables: s.selectables,
		alias:       alias,
	}
}

func (s SubQuery) Col(fieldName string) Column {
	return Column{
		tableRef:  s,
		fieldName: fieldName,
	}
}

func (s SubQuery) InnerJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeInner,
		left:  s,
		right: right,
	}
}

func (s SubQuery) LeftJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeLeft,
		left:  s,
		right: right,
	}
}

func (s SubQuery) RightJoin(right TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTypeRight,
		left:  s,
		right: right,
	}
}
