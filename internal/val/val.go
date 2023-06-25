package val

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/model"
)

type Val interface {
	WriteCols(rows *sql.Rows) error
}

type Creator func(m *model.Model, v any) Val
