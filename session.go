package easyorm

import (
	"context"
	"database/sql"
)

var (
	_ session = (*DB)(nil)
	_ session = (*Tx)(nil)
)

type session interface {
	getCore() *Core
	queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}
