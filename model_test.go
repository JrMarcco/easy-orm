package orm

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
		wantRes *model
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
			wantRes: &model{
				tbName: "parse_model_arg",
				fds: map[string]*field{
					"ID": {
						colName: "id",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name",
					},
				},
			},
		}, {
			name: "basic pointer",
			arg:  &parseModelArg{},
			wantRes: &model{
				tbName: "parse_model_arg",
				fds: map[string]*field{
					"ID": {
						colName: "id",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name",
					},
				},
			},
		}, {
			name: "basic tag",
			arg: func() any {
				type Demo struct {
					ID uint64 `orm:"column=id_suffix"`
				}
				return Demo{}
			}(),
			wantRes: &model{
				tbName: "demo",
				fds: map[string]*field{
					"ID": {
						colName: "id_suffix",
					},
				},
			},
		}, {
			name: "tag with space",
			arg: func() any {
				type Demo struct {
					ID uint64 `orm:" column  = id_suffix"`
				}
				return Demo{}
			}(),
			wantRes: &model{
				tbName: "demo",
				fds: map[string]*field{
					"ID": {
						colName: "id_suffix",
					},
				},
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
			wantRes: &model{
				tbName: "custom_name",
				fds: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
			},
		}, {
			name: "custom name ptr",
			arg:  &customTbNamePtr{},
			wantRes: &model{
				tbName: "custom_name",
				fds: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
			},
		}, {
			name: "empty custom name",
			arg:  emptyCustomTbName{},
			wantRes: &model{
				tbName: "empty_custom_tb_name",
				fds: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
			},
		},
	}

	r := &registry{
		models: make(map[reflect.Type]*model, 64),
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Get(tc.arg)

			assert.Equal(t, tc.wantErr, err)

			if err == nil {
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
		wantRes *model
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
			wantRes: &model{
				tbName: "parse_model_arg",
				fds: map[string]*field{
					"ID": {
						colName: "id",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name",
					},
				},
			},
		}, {
			name: "basic pointer",
			arg:  &parseModelArg{},
			wantRes: &model{
				tbName: "parse_model_arg",
				fds: map[string]*field{
					"ID": {
						colName: "id",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name",
					},
				},
			},
		},
	}

	r := &registry{
		models: make(map[reflect.Type]*model, 8),
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.parseModel(tc.arg)

			assert.Equal(t, tc.wantErr, err)

			if err == nil {
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
		opts    []ModelOpt
		wantErr error
		wantRes *model
	}{
		{
			name:   "basis",
			entity: parseModelArg{},
			wantRes: &model{
				tbName: "parse_model_arg",
				fds: map[string]*field{
					"ID": {
						colName: "id",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name",
					},
				},
			},
		}, {
			name:   "with table name opt",
			entity: parseModelArg{},
			opts: []ModelOpt{
				ModelWithTbName("tb_name"),
			},
			wantRes: &model{
				tbName: "tb_name",
				fds: map[string]*field{
					"ID": {
						colName: "id",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name",
					},
				},
			},
		}, {
			name:   "with empty table name opt",
			entity: parseModelArg{},
			opts: []ModelOpt{
				ModelWithTbName(""),
			},
			wantErr: errs.EmptyTbNameErr,
		}, {
			name:   "with single column name opt",
			entity: parseModelArg{},
			opts: []ModelOpt{
				ModelWithColumName("ID", "id_suffix"),
			},
			wantRes: &model{
				tbName: "parse_model_arg",
				fds: map[string]*field{
					"ID": {
						colName: "id_suffix",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name",
					},
				},
			},
		}, {
			name:   "with multi column name opt",
			entity: parseModelArg{},
			opts: []ModelOpt{
				ModelWithColumName("ID", "id_suffix"),
				ModelWithColumName("NickName", "nick_name_suffix"),
			},
			wantRes: &model{
				tbName: "parse_model_arg",
				fds: map[string]*field{
					"ID": {
						colName: "id_suffix",
					},
					"Name": {
						colName: "name",
					},
					"NickName": {
						colName: "nick_name_suffix",
					},
				},
			},
		}, {
			name:   "with invalid column field opt",
			entity: parseModelArg{},
			opts: []ModelOpt{
				ModelWithColumName("Invalid", "column_name"),
			},
			wantErr: errs.InvalidColumnFdErr("Invalid"),
		}, {
			name:   "with empty column name opt",
			entity: parseModelArg{},
			opts: []ModelOpt{
				ModelWithColumName("Name", ""),
			},
			wantErr: errs.EmptyColNameErr,
		},
	}

	r := newRegistry()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.Register(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}
