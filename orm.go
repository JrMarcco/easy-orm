package easyorm

import (
	"context"
	"database/sql"
)

var (
	_ orm = (*DB)(nil)
	_ orm = (*Tx)(nil)
)

type orm interface {
	getCore() *core
	queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error)
	execContext(ctx context.Context, sql string, args ...any) (sql.Result, error)
}
