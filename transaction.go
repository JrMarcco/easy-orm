package easyorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

var _ driver.Tx = (*Tx)(nil)

type Tx struct {
	sqlTx *sql.Tx
}

func (t *Tx) getCore() *Core {
	//TODO implement me
	panic("implement me")
}

func (t *Tx) queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tx) execContext(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tx) Commit() error {
	//TODO implement me
	panic("implement me")
}

func (t *Tx) Rollback() error {
	//TODO implement me
	panic("implement me")
}
