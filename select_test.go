package easyorm

import (
	"database/sql"
	"testing"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testModel struct {
	Id       uint64
	Name     string
	Age      int8
	NickName *sql.NullString
}

func TestSelector_Build(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, DBWithDialect(MySQLDialect))
	require.NoError(t, err)

	tcs := []struct {
		name          string
		selector      StatementBuilder
		wantStatement *Statement
		wantErr       error
	}{
		{
			name:     "basic select",
			selector: NewSelector[testModel](db),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model`;",
			},
			wantErr: nil,
		}, {
			name:     "select with where",
			selector: NewSelector[testModel](db).Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "select with where not",
			selector: NewSelector[testModel](db).Where(Col("Id").Eq(1).Not()),
			wantStatement: &Statement{
				Sql: "SELECT * FROM `test_model` WHERE NOT (`id` = ?);",
				Args: []any{
					1,
				},
			},
		}, {
			name: "select with where and",
			selector: NewSelector[testModel](db).Where(
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
			selector: NewSelector[testModel](db).Where(
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
		}, {
			name: "select with invalid column",
			selector: NewSelector[testModel](db).Where(
				Col("InvalidColumn").Eq(1),
			),
			wantErr: errs.ErrInvalidField("InvalidColumn"),
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
