package orm

import (
    "context"
    "database/sql"
    "github.com/jrmarcco/easy-orm/internal/val"
    "github.com/jrmarcco/easy-orm/model"
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
