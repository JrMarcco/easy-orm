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

	// SeqFds 用于按顺序记录字段
	// go 的 map 遍历是无序的，
	// 在 insert 这种需要保证字段顺序的情况下则需要通过 SeqFds 来处理。
	SeqFds []*Field
}

type Opt func(m *Model) error

func WithTbName(tb string) Opt {
	return func(m *Model) error {
		if tb == "" {
			return errs.ErrEmptyTbName
		}

		m.Tb = tb
		return nil
	}
}

func WithColName(fdName string, colName string) Opt {
	return func(m *Model) error {

		if colName == "" {
			return errs.ErrEmptyColName

		}

		fd, ok := m.Fds[fdName]
		if !ok {
			return errs.ErrInvalidColumnFd(fdName)
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
