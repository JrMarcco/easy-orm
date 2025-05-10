package easyorm

import (
	"context"
	"database/sql"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/internal/value"
	"github.com/JrMarcco/easy-orm/model"
)

// DB a decorator for sql.DB
type DB struct {
	*Core
	sqlDB *sql.DB
}

func (db *DB) getCore() *Core {
	return db.Core
}

func (db *DB) queryContext(ctx context.Context, sql string, args ...any) (*sql.Rows, error) {
	return db.sqlDB.QueryContext(ctx, sql, args...)
}

func (db *DB) execContext(ctx context.Context, sql string, args ...any) (sql.Result, error) {
	return db.sqlDB.ExecContext(ctx, sql, args...)
}

func (db *DB) beginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	sqlTx, err := db.sqlDB.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Tx{
		sqlTx: sqlTx,
	}, nil
}

func (db *DB) DoTx(ctx context.Context, bizFunc func(ctx context.Context, tx *Tx) error, opts *sql.TxOptions) (err error) {
	tx, err := db.beginTx(ctx, opts)
	if err != nil {
		return err
	}

	panicked := true
	defer func() {
		if panicked || err != nil {
			rbErr := tx.Rollback()
			err = errs.ErrRollback(err, rbErr, panicked)
			return
		}

		err = tx.Commit()
	}()

	err = bizFunc(ctx, tx)
	panicked = false

	return err
}

type DBOpt func(db *DB)

func DBWithDialect(dialect mysql) DBOpt {
	return func(db *DB) {
		db.dialect = dialect
	}
}

func DBWithRegistry(registry model.Registry) DBOpt {
	return func(db *DB) {
		db.registry = registry
	}
}

func DBWithValueResolver(resolverCreator value.ResolverCreator) DBOpt {
	return func(db *DB) {
		db.resolverCreator = resolverCreator
	}
}

func Open(driverName string, dsn string, opts ...DBOpt) (*DB, error) {
	sqlDB, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	return OpenDB(sqlDB, opts...)
}

func OpenDB(sqlDB *sql.DB, opts ...DBOpt) (*DB, error) {
	core := &Core{
		registry:        model.NewRegistry(),
		resolverCreator: value.NewUnsafeResolver,
	}

	db := &DB{
		Core:  core,
		sqlDB: sqlDB,
	}

	for _, opt := range opts {
		opt(db)
	}
	return db, nil
}
