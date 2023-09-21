package orm

import "database/sql"

// Result 结果集封装
type Result struct {
	res sql.Result
	err error
}

func (r Result) Err() error {
	return r.err
}

func (r Result) LastInsertId() int64 {
	id, err := r.res.LastInsertId()
	if err != nil {
		r.err = err
		return 0
	}

	return id
}

func (r Result) RowsAffected() int64 {
	rowsAffected, err := r.res.RowsAffected()
	if err != nil {
		r.err = err
		return 0
	}

	return rowsAffected
}
