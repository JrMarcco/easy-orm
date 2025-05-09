package value

import (
	"database/sql"
	"reflect"
	"unsafe"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/model"
)

var _ ValResolver = (*unsafeResolver)(nil)

type unsafeResolver struct {
	model *model.Model
	addr  unsafe.Pointer
}

func (u unsafeResolver) ReadColumn(fieldName string) (any, error) {
	field, ok := u.model.Fields[fieldName]
	if !ok {
		return nil, errs.ErrInvalidField(fieldName)
	}

	p := unsafe.Pointer(uintptr(u.addr) + field.Offset)

	val := reflect.NewAt(field.Typ, p)
	return val.Elem().Interface(), nil
}

func (u unsafeResolver) WriteColumns(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	values := make([]any, 0, len(columns))

	for _, column := range columns {
		field, ok := u.model.Columns[column]
		if !ok {
			return errs.ErrInvalidColumn(column)
		}

		p := unsafe.Pointer(uintptr(u.addr) + field.Offset)
		if p == nil {
			return errs.ErrInvalidColumn(column)
		}

		val := reflect.NewAt(field.Typ, p)
		values = append(values, val.Interface())
	}

	if err = rows.Scan(values...); err != nil {
		return err
	}
	return nil
}

func NewUnsafeResolver(model *model.Model, v any) ValResolver {
	return unsafeResolver{
		model: model,
		addr:  reflect.ValueOf(v).UnsafePointer(),
	}
}
