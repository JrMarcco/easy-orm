package orm

import (
	"context"
	"database/sql"
)

type Tx struct {
	*Core
	sqlTx *sql.Tx
}

func (t *Tx) getCore() *Core {
	return t.Core
}

func (t *Tx) queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	return t.sqlTx.QueryContext(ctx, sql, args...)
}

func (t *Tx) execContext(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	return t.sqlTx.ExecContext(ctx, sql, args...)
}

func (t *Tx) Commit() error {
	return t.sqlTx.Commit()

}

func (t *Tx) Rollback() error {
	return t.sqlTx.Rollback()
}
