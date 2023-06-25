package model

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"reflect"
)

type TbName interface {
	TbName() string
}

type Model struct {
	Tb   string
	Fds  map[string]*Field // fieldName -> Field
	Cols map[string]*Field // ColName -> Field
}

type Opt func(m *Model) error

func WithTbName(tb string) Opt {
	return func(m *Model) error {
		if tb == "" {
			return errs.EmptyTbNameErr
		}

		m.Tb = tb
		return nil
	}
}

func WithColName(fdName string, colName string) Opt {
	return func(m *Model) error {

		if colName == "" {
			return errs.EmptyColNameErr

		}

		fd, ok := m.Fds[fdName]
		if !ok {
			return errs.InvalidColumnFdErr(fdName)
		}

		delete(m.Cols, fd.ColName)
		m.Cols[colName] = fd

		fd.ColName = colName
		return nil
	}
}

type Field struct {
	Type    reflect.Type
	Name    string
	ColName string
	Offset  uintptr
}
