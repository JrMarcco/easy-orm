package easyorm

import (
	"database/sql"
	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testModel struct {
	Id       uint64
	Name     string
	Age      int8
	NickName *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	tcs := []struct {
		name          string
		selector      StatementBuilder
		wantStatement *Statement
		wantErr       error
	}{
		{
			name:     "basic select",
			selector: NewSelector[testModel](),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model`;",
			},
			wantErr: nil,
		}, {
			name:     "select from table",
			selector: NewSelector[testModel]().From("table_name"),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `table_name`;",
			},
		}, {
			name:     "select with empty table fieldName",
			selector: NewSelector[testModel]().From(""),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model`;",
			},
		}, {
			name:     "select from table with db fieldName",
			selector: NewSelector[testModel]().From("db_name.table_name"),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `db_name`.`table_name`;",
			},
		}, {
			name:     "select from invalid table fieldName",
			selector: NewSelector[testModel]().From("db_name.table_name.sub_table_name"),
			wantErr:  errs.ErrInvalidTableName,
		}, {
			name:     "select with where",
			selector: NewSelector[testModel]().Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "select with where not",
			selector: NewSelector[testModel]().Where(Not(Col("Id").Eq(1))),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model` WHERE NOT (`id` = ?);",
				Args: []any{
					1,
				},
			},
		}, {
			name: "select with where and",
			selector: NewSelector[testModel]().Where(
				Col("Id").Eq(1),
				Col("Name").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model` WHERE (`id` = ?) AND (`name` = ?);",
				Args: []any{
					1, "jrmarcco",
				},
			},
		}, {
			name: "select with where or",
			selector: NewSelector[testModel]().Where(
				Col("Age").Ge(18).
					Or(
						Col("Age").Lt(12),
					),
				Col("NickName").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model` WHERE ((`age` >= ?) OR (`age` < ?)) AND (`nick_name` = ?);",
				Args: []any{
					18, 12, "jrmarcco",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			statement, err := tc.selector.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantStatement, statement)
			}
		})
	}
}
