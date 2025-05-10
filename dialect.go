package easyorm

import "strconv"

var (
	StandardSQL     = standardSQL{}
	PostgresDialect = postgres{}
	MySQLDialect    = mysql{}
)

type Dialect interface {
	quote() byte
	bindArg(b *builder)
}

var _ Dialect = (*standardSQL)(nil)

type standardSQL struct{}

func (s standardSQL) quote() byte {
	return '"'
}

func (s standardSQL) bindArg(b *builder) {
	b.sqlBuffer.WriteByte('?')
}

var _ Dialect = (*postgres)(nil)

type postgres struct {
	standardSQL
}

func (p postgres) bindArg(b *builder) {
	b.sqlBuffer.WriteByte('$')
	b.sqlBuffer.WriteString(strconv.Itoa(len(b.args)))
}

var _ Dialect = (*mysql)(nil)

type mysql struct {
	standardSQL
}

func (m mysql) quote() byte {
	return '`'
}
