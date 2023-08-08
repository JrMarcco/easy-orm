package val

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/model"
	"reflect"
	"unsafe"
)

type unsafeVal struct {
	m    *model.Model
	addr unsafe.Pointer
}

var _ Creator = NewUnsafeValWriter

func NewUnsafeValWriter(m *model.Model, v any) Value {
	return unsafeVal{
		m:    m,
		addr: reflect.ValueOf(v).UnsafePointer(),
	}
}

func (u unsafeVal) ReadCol(fnName string) (any, error) {
	fd, ok := u.m.Fds[fnName]
	if !ok {
		return nil, errs.InvalidColumnFdErr(fnName)
	}

	ptr := unsafe.Pointer(uintptr(u.addr) + fd.Offset)

	// 注意这里 val := reflect.NewAt(fd.Type, ptr)
	// 创建出来的是 fd.fdType 类型的指针。
	val := reflect.NewAt(fd.Type, ptr)
	return val.Elem().Interface(), nil
}

func (u unsafeVal) WriteCols(rows *sql.Rows) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	vals := make([]any, 0, len(cols))

	for _, col := range cols {
		fd, ok := u.m.Cols[col]
		if !ok {
			return errs.InvalidColumnErr(col)
		}

		ptr := unsafe.Pointer(uintptr(u.addr) + fd.Offset)
		if ptr == nil {
			return errs.InvalidColumnErr(col)
		}

		// 注意这里 val := reflect.NewAt(fd.Type, ptr)
		// 创建出来的是 fd.fdType 类型的指针。
		val := reflect.NewAt(fd.Type, ptr)
		vals = append(vals, val.Interface())

	}

	if err = rows.Scan(vals...); err != nil {
		return err
	}

	return nil
}
