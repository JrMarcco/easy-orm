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
				fds: map[string]field{
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
				fds: map[string]field{
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
				fds: map[string]field{
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
				fds: map[string]field{
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
		},
	}

	r := newRegistry()
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.parseModel(tc.arg)

			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
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
				fds: map[string]field{
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
				fds: map[string]field{
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

	r := newRegistry()
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := r.getModel(tc.arg)

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

type parseModelArg struct {
	ID       int64
	Name     string
	NickName string
}
