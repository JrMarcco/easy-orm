package model

import (
	"reflect"
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

type Opt func(*Model)

func WithTableNameOpt(tableName string) Opt {
	return func(m *Model) {
		m.TableName = tableName
	}
}

func WithColumnOpt(fieldName, ColumnName string) Opt {
	return func(m *Model) {
	}
}

type Field struct {
	Typ        reflect.Type
	FiledName  string
	ColumnName string
}
