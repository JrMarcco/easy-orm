package orm

import "github.com/jrmarcco/easy-orm/internal/errs"

var (
	StandardSQL  = standardSQL{}
	MySqlDialect = mysql{}
)

type Dialect interface {
	quote() byte
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

func (s standardSQL) onConflict(b *builder, onConflict *OnConflict) error {
	b.sb.WriteString(" ON CONFLICT (")

	for idx, conflict := range onConflict.conflicts {
		if idx > 0 {
			b.sb.WriteByte(',')
		}
		fd, ok := b.model.Fds[conflict]
		if !ok {
			return errs.InvalidColumnFdErr(conflict)
		}
		b.writeQuote(fd.ColName)
	}

	b.sb.WriteString(") DO UPDATE SET ")

	for idx, assign := range onConflict.assigns {
		if idx > 0 {
			b.sb.WriteByte(',')
		}

		switch typ := assign.(type) {
		case Assignment:
			if err := b.buildAssign(typ); err != nil {
				return err
			}
		case Column:
			fd, ok := b.model.Fds[typ.fdName]
			if !ok {
				return errs.InvalidColumnFdErr(typ.fdName)
			}

			b.writeQuote(fd.ColName)
			b.sb.WriteString("=EXCLUDED.")

			ufd, ok := b.model.Fds[typ.ufdName]
			if !ok {
				return errs.InvalidColumnFdErr(typ.ufdName)
			}

			b.writeQuote(ufd.ColName)
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
			if err := b.buildAssign(typ); err != nil {
				return err
			}
		case Column:
			typ.alias = ""
			if err := b.buildCol(typ); err != nil {
				return err
			}

			b.sb.WriteString("=VALUES(")

			ufd, ok := b.model.Fds[typ.ufdName]
			if !ok {
				return errs.InvalidColumnFdErr(typ.ufdName)
			}

			b.writeQuote(ufd.ColName)
			b.sb.WriteByte(')')
		default:
			return errs.InvalidAssignmentErr
		}
	}
	return nil
}
