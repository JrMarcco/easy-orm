package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"strconv"
)

var (
	StandardSQL     = standardSQL{}
	MySqlDialect    = mysql{}
	PostgresDialect = postgres{}
)

type Dialect interface {
	quote() byte
	bindArg(b *builder)
	onConflict(b *builder, onConflict *OnConflict) error
}

type OnConflict struct {
	conflicts []string
	assigns   []Assignable
}

type standardSQL struct {
}

var _ Dialect = new(standardSQL)

func (s standardSQL) quote() byte {
	return '"'
}

func (s standardSQL) bindArg(b *builder) {
	b.sb.WriteByte('?')
}

func (s standardSQL) onConflict(b *builder, onConflict *OnConflict) error {
	b.sb.WriteString(" ON CONFLICT (")

	for idx, conflict := range onConflict.conflicts {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		if err := b.writeField(conflict); err != nil {
			return err
		}
	}

	b.sb.WriteString(") DO UPDATE SET ")

	for idx, assign := range onConflict.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}

		switch typ := assign.(type) {
		case Assignment:
			if err := b.writeField(typ.fdName); err != nil {
				return err
			}

			b.addArg(typ.val)
			b.sb.WriteByte('=')
			s.bindArg(b)
		case Column:
			typ.alias = ""
			if err := b.buildCol(typ); err != nil {
				return err
			}

			b.sb.WriteString("=EXCLUDED.")

			if err := b.writeField(typ.ufdName); err != nil {
				return err
			}
		}
	}
	return nil
}

type mysql struct {
	standardSQL
}

var _ Dialect = new(mysql)

func (m mysql) quote() byte {
	return '`'
}

func (m mysql) onConflict(b *builder, onConflict *OnConflict) error {

	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")

	for idx, assign := range onConflict.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}

		switch typ := assign.(type) {
		case Assignment:
			if err := b.writeField(typ.fdName); err != nil {
				return err
			}

			b.addArg(typ.val)
			b.sb.WriteByte('=')
			m.bindArg(b)
		case Column:
			typ.alias = ""
			if err := b.buildCol(typ); err != nil {
				return err
			}

			b.sb.WriteString("=VALUES(")

			if err := b.writeField(typ.ufdName); err != nil {
				return err
			}

			b.sb.WriteByte(')')
		default:
			return errs.ErrInvalidAssignment
		}
	}
	return nil
}

type postgres struct {
	standardSQL
}

func (p postgres) bindArg(b *builder) {
	b.sb.WriteByte('$')
	b.sb.WriteString(strconv.Itoa(len(b.args)))
}
