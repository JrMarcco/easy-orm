package orm

import "strings"

type Deletor[T any] struct {
	*builder
	conds []condition
	db    *DB
}

func NewDeletor[T any](db *DB) *Deletor[T] {
	return &Deletor[T]{
		builder: newBuilder(),
		db:      db,
	}
}

func (d *Deletor[T]) From(tbName string) *Deletor[T] {
	d.tbName = tbName
	return d
}

func (d *Deletor[T]) Where(predicates ...Predicate) *Deletor[T] {
	if d.conds == nil {
		d.conds = make([]condition, 0, 2)
	}

	if len(predicates) > 0 {
		d.conds = append(d.conds, newCond(condTypWhere, predicates))
	}
	return d
}

func (d *Deletor[T]) Build() (*Statement, error) {

	var err error
	if d.model, err = d.db.registry.Get(new(T)); err != nil {
		return nil, err
	}

	d.sb = &strings.Builder{}
	d.sb.WriteString("DELETE FROM ")

	if d.tbName == "" {
		d.sb.WriteByte('`')
		d.sb.WriteString(d.model.Tb)
		d.sb.WriteByte('`')
	} else {

		segs := strings.SplitN(d.tbName, ".", 2)

		d.sb.WriteByte('`')
		d.sb.WriteString(segs[0])
		d.sb.WriteByte('`')

		if len(segs) > 1 {
			d.sb.WriteByte('.')
			d.sb.WriteByte('`')
			d.sb.WriteString(segs[1])
			d.sb.WriteByte('`')
		}

	}

	if len(d.conds) > 0 {
		for _, cond := range d.conds {
			d.sb.WriteString(string(cond.typ))

			if err := d.buildExpr(cond.rootExpr); err != nil {
				return nil, err
			}
		}
	}

	d.sb.WriteByte(';')

	return &Statement{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}
