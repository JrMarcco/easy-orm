package model

import (
	"reflect"
	"strings"

	"github.com/JrMarcco/easy-orm/internal/errs"
)

type Registry interface {
	GetModel(entity any) (*Model, error)
	RegisterModel(entity any, opts ...Opt) (*Model, error)
}

type Model struct {
	TableName string
	Fields    map[string]*Field // fieldName -> Field
	Columns   map[string]*Field // ColumnName -> Field
}

type Opt func(*Model) error

func WithTableOpt(tableName string) Opt {
	return func(m *Model) error {
		if tableName == "" {
			return errs.ErrInvalidTable(tableName)
		}

		if segments := strings.Split(tableName, "."); len(segments) > 2 {
			return errs.ErrInvalidTable(tableName)
		}

		m.TableName = tableName
		return nil
	}
}

func WithColumnOpt(fieldName, columnName string) Opt {
	return func(m *Model) error {
		if columnName == "" {
			return errs.ErrInvalidColumn(columnName)
		}

		field, ok := m.Fields[fieldName]
		if !ok {
			return errs.ErrInvalidField(fieldName)
		}

		delete(m.Columns, field.ColumnName)
		m.Columns[columnName] = field

		field.ColumnName = columnName
		return nil
	}
}

type Field struct {
	Typ        reflect.Type
	FiledName  string
	ColumnName string
	Offset     uintptr
}
