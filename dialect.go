package easyorm

import (
	"strconv"

	"github.com/JrMarcco/easy-orm/internal/errs"
)

var (
	StandardSQL     = standardSQL{}
	PostgresDialect = postgres{}
	MySQLDialect    = mysql{}
)

type Dialect interface {
	quote() byte
	bindArg(b *builder)
	onConflict(b *builder, conflict *Conflict) error
}

type Conflict struct {
	conflicts []string // conflict fields
	assigns   []Assignable
}

var _ Dialect = (*standardSQL)(nil)

type standardSQL struct{}

func (s standardSQL) quote() byte {
	return '"'
}

func (s standardSQL) bindArg(b *builder) {
	b.sqlBuffer.WriteByte('?')
}

func (s standardSQL) onConflict(_ *builder, _ *Conflict) error {
	return errs.ErrUnsupportedOnConflict
}

var _ Dialect = (*postgres)(nil)

type postgres struct {
	standardSQL
}

func (p postgres) bindArg(b *builder) {
	b.sqlBuffer.WriteByte('$')
	b.sqlBuffer.WriteString(strconv.Itoa(len(b.args)))
}

func (p postgres) onConflict(b *builder, conflict *Conflict) error {
	b.sqlBuffer.WriteString(" ON CONFLICT (")

	for index, c := range conflict.conflicts {
		if index > 0 {
			b.sqlBuffer.WriteString(", ")
		}
		if err := b.writeField(c); err != nil {
			return err
		}
	}

	b.sqlBuffer.WriteByte(')')

	if len(conflict.assigns) == 0 {
		b.sqlBuffer.WriteString(" DO NOTHING")
		return nil
	}

	b.sqlBuffer.WriteString(" DO UPDATE SET ")
	for index, assign := range conflict.assigns {
		if index > 0 {
			b.sqlBuffer.WriteString(", ")
		}

		switch assignTyp := assign.(type) {
		case Assignment:
			if err := b.writeField(assignTyp.filedName); err != nil {
				return err
			}

			b.addArgs(assignTyp.value)
			b.sqlBuffer.WriteString(" = ")
			b.dialect.bindArg(b)
		case Column:
			assignTyp.alias = ""
			if err := b.buildColumn(assignTyp.tableRef, assignTyp.fieldName); err != nil {
				return err
			}

			b.sqlBuffer.WriteString(" = EXCLUDED.")
			if err := b.writeField(assignTyp.fieldName); err != nil {
				return err
			}
		default:
			return errs.ErrInvalidAssignable
		}
	}
	return nil
}

var _ Dialect = (*mysql)(nil)

type mysql struct {
	standardSQL
}

func (m mysql) quote() byte {
	return '`'
}

func (m mysql) onConflict(b *builder, conflict *Conflict) error {
	b.sqlBuffer.WriteString(" ON DUPLICATE KEY UPDATE ")

	for index, assign := range conflict.assigns {
		if index > 0 {
			b.sqlBuffer.WriteString(", ")
		}

		switch assignTyp := assign.(type) {
		case Assignment:
			if err := b.writeField(assignTyp.filedName); err != nil {
				return err
			}

			b.addArgs(assignTyp.value)
			b.sqlBuffer.WriteString(" = ")
			b.dialect.bindArg(b)
		case Column:
			assignTyp.alias = ""
			if err := b.buildColumn(assignTyp.tableRef, assignTyp.fieldName); err != nil {
				return err
			}

			b.sqlBuffer.WriteString(" = VALUES(")
			if err := b.writeField(assignTyp.fieldName); err != nil {
				return err
			}
			b.sqlBuffer.WriteByte(')')
		default:
			return errs.ErrInvalidAssignable
		}
	}
	return nil
}
