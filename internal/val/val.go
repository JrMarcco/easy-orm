package val

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/model"
)

type Value interface {
	// ReadCol 读取单个字段值
	ReadCol(fdName string) (any, error)
	// WriteCols 写入整行字段值
	WriteCols(rows *sql.Rows) error
}

type Creator func(m *model.Model, v any) Value
