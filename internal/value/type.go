package value

import (
	"database/sql"
)

type ValResolver interface {
	ReadColumn(fieldName string) (any, error)
	WriteColumns(rows *sql.Rows) error
}
