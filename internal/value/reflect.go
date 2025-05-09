package value

import (
	"database/sql"
	"reflect"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/model"
)

var _ ValResolver = (*reflectResolver)(nil)

type reflectResolver struct {
	model *model.Model
	val   reflect.Value
}

func (r reflectResolver) ReadColumn(fieldName string) (any, error) {
	return r.val.FieldByName(fieldName).Interface(), nil
}

func (r reflectResolver) WriteColumns(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := make([]any, 0, len(columns))
	valueElements := make([]reflect.Value, 0, len(columns))

	for _, column := range columns {
		field, ok := r.model.Columns[column]
		if !ok {
			return errs.ErrInvalidColumn(column)
		}

		val := reflect.New(field.Typ)
		values = append(values, val.Interface())
		valueElements = append(valueElements, val.Elem())
	}

	if err = rows.Scan(values...); err != nil {
		return err
	}

	for i, column := range columns {
		field, _ := r.model.Columns[column]

		r.val.FieldByName(field.FiledName).Set(valueElements[i])
	}
	return nil
}

func NewReflectResolver(model *model.Model, v any) ValResolver {
	return reflectResolver{
		model: model,
		val:   reflect.ValueOf(v).Elem(),
	}
}
