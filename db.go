package orm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/internal/val"
	"github.com/jrmarcco/easy-orm/model"
	"log"
	"time"
)

type DB struct {
	*Core
	sqlDB *sql.DB
}

type DBOpt func(db *DB)

func DDWithValCreator(creator val.Creator) DBOpt {
	return func(db *DB) {
		db.creator = creator
	}
}

func DBWithDialect(dialect Dialect) DBOpt {
	return func(db *DB) {
		db.dialect = dialect
	}
}

func DBWithMdls(mdls ...Middleware) DBOpt {
	return func(db *DB) {
		db.mdls = mdls
	}
}

func Open(driver string, dsn string, opts ...DBOpt) (*DB, error) {
	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return OpenDB(sqlDB, opts...)
}

func OpenDB(sqlDB *sql.DB, opts ...DBOpt) (*DB, error) {

	core := &Core{
		registry: model.NewRegistry(),
		creator:  val.NewUnsafeValWriter,
		dialect:  StandardSQL,
	}

	db := &DB{
		sqlDB: sqlDB,
		Core:  core,
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}

func (d *DB) getCore() *Core {
	return d.Core
}

func (d *DB) Wait() error {
	err := d.sqlDB.Ping()

	for errors.Is(err, driver.ErrBadConn) {
		log.Println("waiting for db start ...")
		err = d.sqlDB.Ping()

		time.Sleep(time.Second)
	}

	return err
}

func (d *DB) queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	return d.sqlDB.QueryContext(ctx, sql, args...)
}

func (d *DB) execContext(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	return d.sqlDB.ExecContext(ctx, sql, args...)
}

func (d *DB) beginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	sqlTx, err := d.sqlDB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		Core:  d.Core,
		sqlTx: sqlTx,
	}, nil
}

func (d *DB) DoTransaction(ctx context.Context, bizFunc func(ctx context.Context, tx *Tx), opts *sql.TxOptions) (err error) {

	tx, err := d.beginTx(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		// do rollback
		if panicked || err != nil {
			rollbackErr := tx.Rollback()
			err = errs.ErrRollback(err, rollbackErr, panicked)
			return
		}
		// do commit
		err = tx.Commit()
	}()

	bizFunc(ctx, tx)
	panicked = false

	return err
}
