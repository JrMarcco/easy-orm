package testdata

import "database/sql"

type Model struct {
	Id       uint64
	Age      *int32
	Username string
	Address  *sql.NullString
}

type SubModel struct {
	Id      uint64
	Name    string
	Email   *sql.NullString
	Balance float64
}
