package model

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type basicStruct struct {
	Id       uint64
	Age      int8
	Name     string
	NickName *sql.NullString
	IDCardNo string
}

type withTagStruct struct {
	Id   uint64
	Name string `orm:"column=user_name"`
}

type withInvalidTagStruct1 struct {
	Id   uint64
	Name string `orm:"column="`
}
type withInvalidTagStruct2 struct {
	Id   uint64
	Name string `orm:"=user_name"`
}

type withInvalidTagStruct3 struct {
	Id   uint64
	Name string `orm:"column-user_name"`
}

func TestModelRegistry_RegisterModel(t *testing.T) {
	r := NewRegistry()
	tcs := []struct {
		name      string
		entity    any
		opts      []Opt
		wantModel *Model
		wantErr   error
	}{
		{
			name:   "basic struct",
			entity: basicStruct{},
			wantModel: &Model{
				TableName: "basic_struct",
				SeqFields: []*Field{
					{
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					}, {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					}, {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"NickName": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"IDCardNo": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"nick_name": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"id_card_no": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
			},
			wantErr: nil,
		}, {
			name:   "pointer to struct",
			entity: &basicStruct{},
			wantModel: &Model{
				TableName: "basic_struct",
				SeqFields: []*Field{
					{
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					}, {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					}, {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"NickName": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"IDCardNo": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"nick_name": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"id_card_no": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
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
		}, {
			name:   "struct with model opts",
			entity: basicStruct{},
			opts: []Opt{
				WithTableOpt("t_table"),
				WithColumnOpt("Name", "user_name"),
			},
			wantModel: &Model{
				TableName: "t_table",
				SeqFields: []*Field{
					{
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					}, {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     16,
					}, {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     16,
					},
					"NickName": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"IDCardNo": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"user_name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     16,
					},
					"nick_name": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"id_card_no": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
			},
		}, {
			name:   "pointer with model opts",
			entity: &basicStruct{},
			opts: []Opt{
				WithTableOpt("t_table_name"),
				WithColumnOpt("IDCardNo", "card_no"),
			},
			wantModel: &Model{
				TableName: "t_table_name",
				SeqFields: []*Field{
					{
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					}, {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					}, {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "card_no",
						Offset:     40,
					},
				},
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"NickName": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"IDCardNo": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "card_no",
						Offset:     40,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"nick_name": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"card_no": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "card_no",
						Offset:     40,
					},
				},
			},
		}, {
			name:   "pointer with invalid table opt",
			entity: &basicStruct{},
			opts: []Opt{
				WithTableOpt("t.table.name"),
			},
			wantErr: errs.ErrInvalidTable("t.table.name"),
		}, {
			name:   "pointer with unknown field opt",
			entity: &basicStruct{},
			opts: []Opt{
				WithColumnOpt("UnknownFiled", "unknown_field"),
			},
			wantErr: errs.ErrInvalidField("UnknownFiled"),
		}, {
			name:   "struct with tag",
			entity: withTagStruct{},
			wantModel: &Model{
				TableName: "with_tag_struct",
				SeqFields: []*Field{
					{
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					}, {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     8,
					},
				},
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     8,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"user_name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     8,
					},
				},
			},
		}, {
			name:    "struct with invalid tag 1",
			entity:  withInvalidTagStruct1{},
			wantErr: errs.ErrInvalidTag("column="),
		}, {
			name:    "struct with invalid tag 2",
			entity:  withInvalidTagStruct2{},
			wantErr: errs.ErrInvalidTag("=user_name"),
		}, {
			name:    "struct with invalid tag 3",
			entity:  withInvalidTagStruct3{},
			wantErr: errs.ErrInvalidTag("column-user_name"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			var model *Model
			if len(tc.opts) > 0 {
				model, err = r.RegisterModel(tc.entity, tc.opts...)
			} else {
				model, err = r.RegisterModel(tc.entity)
			}

			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantModel, model)
			}
		})
	}
}

func TestModelRegistry_GetModel(t *testing.T) {
	r := &modelRegistry{
		models: map[reflect.Type]*Model{},
	}
	var err error
	_, err = r.RegisterModel(basicStruct{})
	require.NoError(t, err)
	_, err = r.RegisterModel(&withTagStruct{}, WithTableOpt("t_table"))
	require.NoError(t, err)

	tcs := []struct {
		name      string
		entity    any
		wantModel *Model
		wantErr   error
	}{
		{
			name:   "struct",
			entity: basicStruct{},
			wantModel: &Model{
				TableName: "basic_struct",
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"NickName": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"IDCardNo": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"nick_name": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"id_card_no": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
			},
			wantErr: nil,
		}, {
			name:   "pointer",
			entity: &basicStruct{},
			wantModel: &Model{
				TableName: "basic_struct",
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"NickName": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"IDCardNo": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"age": {
						Typ:        reflect.TypeOf(int8(0)),
						FiledName:  "Age",
						ColumnName: "age",
						Offset:     8,
					},
					"name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "name",
						Offset:     16,
					},
					"nick_name": {
						Typ:        reflect.TypeOf(&sql.NullString{}),
						FiledName:  "NickName",
						ColumnName: "nick_name",
						Offset:     32,
					},
					"id_card_no": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "IDCardNo",
						ColumnName: "id_card_no",
						Offset:     40,
					},
				},
			},
			wantErr: nil,
		}, {
			name:   "struct",
			entity: withTagStruct{},
			wantModel: &Model{
				TableName: "with_tag_struct",
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     8,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"user_name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     8,
					},
				},
			},
			wantErr: nil,
		}, {
			name:   "pointer",
			entity: &withTagStruct{},
			wantModel: &Model{
				TableName: "t_table",
				Fields: map[string]*Field{
					"Id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"Name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     8,
					},
				},
				Columns: map[string]*Field{
					"id": {
						Typ:        reflect.TypeOf(uint64(0)),
						FiledName:  "Id",
						ColumnName: "id",
						Offset:     0,
					},
					"user_name": {
						Typ:        reflect.TypeOf(""),
						FiledName:  "Name",
						ColumnName: "user_name",
						Offset:     8,
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			model, err := r.GetModel(tc.entity)
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantModel, model)
			}
		})
	}

	assert.Equal(t, len(r.models), 4)
}
