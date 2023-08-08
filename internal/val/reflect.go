package val

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/model"
	"reflect"
)

type refVal struct {
	m *model.Model
	v reflect.Value
}

var _ Creator = NewRefValWriter

func NewRefValWriter(m *model.Model, v any) Value {
	return refVal{
		m: m,
		v: reflect.ValueOf(v).Elem(),
	}
}

func (r refVal) ReadCol(fdName string) (any, error) {
	return r.v.FieldByName(fdName).Interface(), nil
}

func (r refVal) WriteCols(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	vals := make([]any, 0, len(cols))
	valElems := make([]reflect.Value, 0, len(cols))

	for _, col := range cols {
		fd, ok := r.m.Cols[col]
		if !ok {
			return errs.InvalidColumnErr(col)
		}

		// 注意这里 val := reflect.New(fd.Type)
		// 创建出来的是 fd.fdType 类型的指针。
		val := reflect.New(fd.Type)
		vals = append(vals, val.Interface())
		valElems = append(valElems, val.Elem())
	}

	if err = rows.Scan(vals...); err != nil {
		return err
	}

	for i, col := range cols {
		fd, ok := r.m.Cols[col]
		if !ok {
			return errs.InvalidColumnErr(col)
		}

		r.v.FieldByName(fd.Name).Set(valElems[i])
	}
	return nil
}
