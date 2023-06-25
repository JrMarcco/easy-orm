package model

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestCamelToUnderline(t *testing.T) {
	tcs := []struct {
		name    string
		camel   string
		wantRes string
	}{
		{
			name:    "basic",
			camel:   "firstSecond",
			wantRes: "first_second",
		}, {
			name:    "continuous upper",
			camel:   "ID",
			wantRes: "id",
		}, {
			name:    "had benn underline",
			camel:   "first_second",
			wantRes: "first_second",
		}, {
			name:    "with number",
			camel:   "first1Second2",
			wantRes: "first1_second2",
		}, {
			name:    "continuous upper with number",
			camel:   "DriverCardID",
			wantRes: "driver_card_id",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := camelToUnderline(tc.camel)
			assert.Equal(t, tc.wantRes, res)
		})
	}
}

func TestRegistry_parseModel(t *testing.T) {
	tcs := []struct {
		name    string
		arg     any
		fds     []*Field
		wantRes *Model
		wantErr error
	}{
		{
			name:    "invalid type slice",
			arg:     []int{},
			wantErr: errs.InvalidTypeErr,
		}, {
			name: "multi-level pointer",
			arg: func() **parseModelArg {
				arg := &parseModelArg{}
				return &arg
			},
			wantErr: errs.InvalidTypeErr,
		}, {
			name: "basic struct",
			arg:  parseModelArg{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "parse_model_arg",
			},
		}, {
			name: "basic pointer",
			arg:  &parseModelArg{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "parse_model_arg",
			},
		}, {
			name: "basic tag",
			arg: func() any {
				type Demo struct {
					ID uint64 `orm:"column=id_suffix"`
				}
				return Demo{}
			}(),
			fds: []*Field{
				{
					Type:    reflect.TypeOf(uint64(0)),
					Name:    "ID",
					ColName: "id_suffix",
				},
			},
			wantRes: &Model{
				Tb: "demo",
			},
		}, {
			name: "tag with space",
			arg: func() any {
				type Demo struct {
					ID uint64 `orm:" column  = id_suffix"`
				}
				return Demo{}
			}(),
			fds: []*Field{
				{
					Type:    reflect.TypeOf(uint64(0)),
					Name:    "ID",
					ColName: "id_suffix",
				},
			},
			wantRes: &Model{
				Tb: "demo",
			},
		}, {
			name: "invalid tag content",
			arg: func() any {
				type Demo struct {
					Name string `orm:"column=user_name=s"`
				}
				return Demo{}
			}(),
			wantErr: errs.InvalidTagContentErr("column=user_name=s"),
		}, {
			name: "empty tag key",
			arg: func() any {
				type Demo struct {
					Name string `orm:"  =val"`
				}
				return Demo{}
			}(),
			wantErr: errs.EmptyTagKeyErr("  =val"),
		}, {
			name: "empty tag val",
			arg: func() any {
				type Demo struct {
					Name string `orm:"column=  "`
				}
				return Demo{}
			}(),
			wantErr: errs.EmptyTagValErr("column=  "),
		}, {
			name: "custom name",
			arg:  &customTbName{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  0,
				},
			},
			wantRes: &Model{
				Tb: "custom_name",
			},
		}, {
			name: "custom name ptr",
			arg:  &customTbNamePtr{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  0,
				},
			},
			wantRes: &Model{
				Tb: "custom_name",
			},
		}, {
			name: "empty custom name",
			arg:  emptyCustomTbName{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  0,
				},
			},
			wantRes: &Model{
				Tb: "empty_custom_tb_name",
			},
		},
	}

	r := &registry{
		models: make(map[reflect.Type]*Model, 64),
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Get(tc.arg)

			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				fds := make(map[string]*Field, len(tc.fds))
				cols := make(map[string]*Field, len(tc.fds))

				for _, fd := range tc.fds {
					f := &Field{
						Type:    fd.Type,
						Name:    fd.Name,
						ColName: fd.ColName,
						Offset:  fd.Offset,
					}

					fds[fd.Name] = f
					cols[fd.ColName] = f
				}

				tc.wantRes.Fds = fds
				tc.wantRes.Cols = cols

				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}

type parseModelArg struct {
	ID       int64
	Name     string
	NickName string
}

type customTbName struct {
	Name string
}

func (c customTbName) TbName() string {
	return "custom_name"
}

type customTbNamePtr struct {
	Name string
}

func (c *customTbNamePtr) TbName() string {
	return "custom_name"
}

type emptyCustomTbName struct {
	Name string
}

func (e emptyCustomTbName) TbName() string {
	return ""
}

func TestRegistry_getModel(t *testing.T) {
	tcs := []struct {
		name    string
		arg     any
		fds     []*Field
		wantRes *Model
		wantErr error
	}{
		{
			name:    "invalid type slice",
			arg:     []int{},
			wantErr: errs.InvalidTypeErr,
		}, {
			name: "multi-level pointer",
			arg: func() **parseModelArg {
				arg := &parseModelArg{}
				return &arg
			},
			wantErr: errs.InvalidTypeErr,
		}, {
			name: "basic struct",
			arg:  parseModelArg{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "parse_model_arg",
			},
		}, {
			name: "basic pointer",
			arg:  &parseModelArg{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "parse_model_arg",
			},
		},
	}

	r := &registry{
		models: make(map[reflect.Type]*Model, 8),
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.parseModel(tc.arg)

			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				fds := make(map[string]*Field, len(tc.fds))
				cols := make(map[string]*Field, len(tc.fds))

				for _, fd := range tc.fds {
					f := &Field{
						Type:    fd.Type,
						Name:    fd.Name,
						ColName: fd.ColName,
						Offset:  fd.Offset,
					}

					fds[fd.Name] = f
					cols[fd.ColName] = f
				}

				tc.wantRes.Fds = fds
				tc.wantRes.Cols = cols

				assert.Equal(t, tc.wantRes, res)

				m, ok := r.models[reflect.TypeOf(tc.arg)]
				assert.True(t, ok)

				assert.Equal(t, tc.wantRes, m)
			}
		})
	}
}

func TestRegistry_Register(t *testing.T) {
	tcs := []struct {
		name    string
		entity  any
		opts    []Opt
		fds     []*Field
		wantErr error
		wantRes *Model
	}{
		{
			name:   "basis",
			entity: parseModelArg{},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "parse_model_arg",
			},
		}, {
			name:   "with table name opt",
			entity: parseModelArg{},
			opts: []Opt{
				WithTbName("tb_name"),
			},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "tb_name",
			},
		}, {
			name:   "with empty table name opt",
			entity: parseModelArg{},
			opts: []Opt{
				WithTbName(""),
			},
			wantErr: errs.EmptyTbNameErr,
		}, {
			name:   "with single column name opt",
			entity: parseModelArg{},
			opts: []Opt{
				WithColName("ID", "id_suffix"),
			},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id_suffix",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "parse_model_arg",
			},
		}, {
			name:   "with multi column name opt",
			entity: parseModelArg{},
			opts: []Opt{
				WithColName("ID", "id_suffix"),
				WithColName("NickName", "nick_name_suffix"),
			},
			fds: []*Field{
				{
					Type:    reflect.TypeOf(int64(0)),
					Name:    "ID",
					ColName: "id_suffix",
					Offset:  0,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "Name",
					ColName: "name",
					Offset:  8,
				}, {
					Type:    reflect.TypeOf(""),
					Name:    "NickName",
					ColName: "nick_name_suffix",
					Offset:  24,
				},
			},
			wantRes: &Model{
				Tb: "parse_model_arg",
			},
		}, {
			name:   "with invalid column Field opt",
			entity: parseModelArg{},
			opts: []Opt{
				WithColName("Invalid", "column_name"),
			},
			wantErr: errs.InvalidColumnFdErr("Invalid"),
		}, {
			name:   "with empty column name opt",
			entity: parseModelArg{},
			opts: []Opt{
				WithColName("Name", ""),
			},
			wantErr: errs.EmptyColNameErr,
		},
	}

	r := NewRegistry()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Register(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				fds := make(map[string]*Field, len(tc.fds))
				cols := make(map[string]*Field, len(tc.fds))

				for _, fd := range tc.fds {
					f := &Field{
						Type:    fd.Type,
						Name:    fd.Name,
						ColName: fd.ColName,
						Offset:  fd.Offset,
					}

					fds[fd.Name] = f
					cols[fd.ColName] = f
				}

				tc.wantRes.Fds = fds
				tc.wantRes.Cols = cols

				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}
