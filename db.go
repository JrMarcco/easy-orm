package orm

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/internal/val"
	"github.com/jrmarcco/easy-orm/model"
)

type DB struct {
	registry model.Registry
	sqlDB    *sql.DB
	creator  val.Creator
}

type DBOpt func(db *DB)

func DDWithValCreator(creator val.Creator) DBOpt {
	return func(db *DB) {
		db.creator = creator
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
	db := &DB{
		registry: model.NewRegistry(),
		sqlDB:    sqlDB,
		creator:  val.NewUnsafeValWriter,
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}
