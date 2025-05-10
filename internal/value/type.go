package value

import (
	"database/sql"

	"github.com/JrMarcco/easy-orm/model"
)

type ValResolver interface {
	// ReadColumn read value from SQL columns and set into entity.
	ReadColumn(fieldName string) (any, error)
	// WriteColumns read value from the entity and set into sql.Rows.
	WriteColumns(rows *sql.Rows) error
}

type ResolverCreator func(model *model.Model, v any) ValResolver
