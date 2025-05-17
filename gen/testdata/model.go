package testdata

import sqlx "database/sql"

type Model struct {
	Id       uint64
	Age      *int32
	Username string
	Address  *sqlx.NullString
}

type SubModel struct {
	Id      uint64
	Name    string
	Email   *sqlx.NullString
	Balance float64
}
