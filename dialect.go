package orm

import "github.com/jrmarcco/easy-orm/internal/errs"

var (
	MySqlDialect = mysql{}
)

type Dialect interface {
	quote() byte
	onConflict(b *builder, onConflict *OnConflict) error
}

type OnConflict struct {
	conflicts []Assignable
}

type standardSQL struct {
}

var _ Dialect = new(standardSQL)

func (s standardSQL) quote() byte {
	//TODO implement me
	panic("implement me")
}

func (s standardSQL) onConflict(b *builder, onConflict *OnConflict) error {
	//TODO implement me
	panic("implement me")
}

type mysql struct {
}

var _ Dialect = new(mysql)

func (m mysql) quote() byte {
	return '`'
}

func (m mysql) onConflict(b *builder, onConflict *OnConflict) error {

	b.sb.WriteString(" ON DUPLICATE KEY UPDATE ")

	for idx, assignable := range onConflict.conflicts {
		if idx > 0 {
			b.sb.WriteByte(',')
		}

		switch typ := assignable.(type) {
		case Assignment:
			fd, ok := b.model.Fds[typ.fdName]
			if !ok {
				return errs.InvalidColumnFdErr(typ.fdName)
			}

			b.writeQuote(fd.ColName)
			b.sb.WriteString("=?")

			b.addArg(typ.val)
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
