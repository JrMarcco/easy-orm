package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
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

func TestParseModel(t *testing.T) {
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

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseModel(tc.arg)

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
