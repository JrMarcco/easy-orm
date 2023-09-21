package orm

type SubQ struct {
	builder StatBuilder
	cols    []selectable
	alias   string
	tbRef   TableRef
}

func (s SubQ) expr() {}

func (s SubQ) table() {}

func (s SubQ) tbAlias() string {
	return s.alias
}

func (s SubQ) Join(target TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   JoinTyp,
		left:  s,
		right: target,
	}
}

func (s SubQ) LeftJoin(target TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   LeftJoinTyp,
		left:  s,
		right: target,
	}
}

func (s SubQ) RightJoin(target TableRef) *JoinBuilder {
	return &JoinBuilder{
		typ:   RightJoinTyp,
		left:  s,
		right: target,
	}
}

func (s SubQ) Col(name string) Column {
	return Column{
		tbRef:  s,
		fdName: name,
	}
}
