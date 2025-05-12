package easyorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
)

var _ driver.Tx = (*Tx)(nil)

type Tx struct {
	*core
	sqlTx *sql.Tx
}

func (t *Tx) getCore() *core {
	return t.core
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

func (t *Tx) RollbackIfNotCommit() error {
	if err := t.sqlTx.Rollback(); !errors.Is(err, sql.ErrTxDone) {
		return err
	}
	return nil
}
