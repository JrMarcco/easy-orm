package orm

import "database/sql"

type Result struct {
	res sql.Result
	err error
}

func (r Result) Err() error {
	return r.err
}

func (r Result) LastInsertId() (int64, error) {
	return r.res.LastInsertId()
}

func (r Result) RowsAffected() (int64, error) {
	return r.res.RowsAffected()
}
