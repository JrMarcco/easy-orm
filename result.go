package easyorm

import "database/sql"

// Result sql execute result.
type Result struct {
	res sql.Result
	err error
}

func (r Result) RowsAffected() int64 {
	rows, err := r.res.RowsAffected()
	if err != nil {
		r.err = err
		return 0
	}
	return rows
}

func (r Result) LastInsertId() int64 {
	id, err := r.res.LastInsertId()
	if err != nil {
		r.err = err
		return 0
	}
	return id
}

func (r Result) Err() error {
	return r.err
}
