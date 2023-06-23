package orm

import "database/sql"

type DB struct {
	registry Registry
	sqlDB    *sql.DB
}

type DBOpt func(db *DB)

func Open(driver string, dsn string, opts ...DBOpt) (*DB, error) {

	sqlDB, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return OpenDB(sqlDB, opts...)
}

func OpenDB(sqlDB *sql.DB, opts ...DBOpt) (*DB, error) {
	db := &DB{
		registry: newRegistry(),
		sqlDB:    sqlDB,
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}
