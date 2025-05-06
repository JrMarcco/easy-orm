package model

import (
	"database/sql"
	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type basicStruct struct {
	Id       uint64
	Age      int8
	Name     string
	NickName *sql.NullString
	IDCardNo string
}

func Test_parseModel(t *testing.T) {
	tcs := []struct {
		name      string
		entity    any
		wantModel *Model
		wantErr   error
	}{
		{
			name:   "basic struct",
			entity: basicStruct{},
			wantModel: &Model{
				TableName: "basic_struct",
				Fields: map[string]*Field{
					"Id": {
						FiledName:  "Id",
						ColumnName: "id",
					},
					"Age": {
						FiledName:  "Age",
						ColumnName: "age",
					},
					"Name": {
						FiledName:  "Name",
						ColumnName: "name",
					},
					"NickName": {
						FiledName:  "NickName",
						ColumnName: "nick_name",
					},
					"IDCardNo": {
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
					},
				},
			},
			wantErr: nil,
		}, {
			name:   "pointer to struct",
			entity: &basicStruct{},
			wantModel: &Model{
				TableName: "basic_struct",
				Fields: map[string]*Field{
					"Id": {
						FiledName:  "Id",
						ColumnName: "id",
					},
					"Age": {
						FiledName:  "Age",
						ColumnName: "age",
					},
					"Name": {
						FiledName:  "Name",
						ColumnName: "name",
					},
					"NickName": {
						FiledName:  "NickName",
						ColumnName: "nick_name",
					},
					"IDCardNo": {
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
					},
				},
			},
		}, {
			name: "pointer to pointer",
			entity: func() **basicStruct {
				bs := &basicStruct{}
				return &bs
			}(),
			wantErr: errs.ErrInvalidModelType,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			model, err := parseModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantModel, model)
			}
		})
	}
}
