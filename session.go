package orm

import (
	"context"
	"database/sql"
)

var (
	_ Session = &DB{}
	_ Session = &Tx{}
)

type Session interface {
	getCore() *Core
	queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}
