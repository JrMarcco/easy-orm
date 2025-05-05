package easyorm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type model struct {
	Id       uint64
	Name     string
	Age      int8
	Nickname *sql.NullString
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
			selector: NewSelector[model](),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `model`;",
			},
			wantErr: nil,
		}, {
			name:     "select from table",
			selector: NewSelector[model]().From("table_name"),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `table_name`;",
			},
		}, {
			name:     "select with empty table name",
			selector: NewSelector[model]().From(""),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `model`;",
			},
		}, {
			name:     "select from table with db name",
			selector: NewSelector[model]().From("db_name.table_name"),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `db_name`.`table_name`;",
			},
		}, {
			name:     "select from invalid table name",
			selector: NewSelector[model]().From("db_name.table_name.sub_table_name"),
			wantErr:  errInvalidTableName,
		}, {
			name:     "select with where",
			selector: NewSelector[model]().Where(Col("id").Eq(1)),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "select with where not",
			selector: NewSelector[model]().Where(Not(Col("id").Eq(1))),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `model` WHERE NOT (`id` = ?);",
				Args: []any{
					1,
				},
			},
		}, {
			name: "select with where and",
			selector: NewSelector[model]().Where(
				Col("id").Eq(1),
				Col("name").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `model` WHERE (`id` = ?) AND (`name` = ?);",
				Args: []any{
					1, "jrmarcco",
				},
			},
		}, {
			name: "select with where or",
			selector: NewSelector[model]().Where(
				Col("age").Ge(18).
					Or(
						Col("age").Lt(12),
					),
				Col("nick_name").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `model` WHERE ((`age` >= ?) OR (`age` < ?)) AND (`nick_name` = ?);",
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
