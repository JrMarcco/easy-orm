package easyorm

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/internal/value"
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
		}, {
			name: "with basic selectable",
			selector: NewSelector[testModel](db).
				Select(
					Col("Id"),
					Col("Name"),
					Col("Age"),
				).
				Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				Sql: "SELECT `id`, `name`, `age` FROM `test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name: "with alias selectable",
			selector: NewSelector[testModel](db).
				Select(
					Col("Id"),
					Col("Name").As("user_name"),
					Col("Age"),
				).
				Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				Sql: "SELECT `id`, `name` AS user_name, `age` FROM `test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "with aggregate",
			selector: NewSelector[testModel](db).Select(Count("Id")),
			wantStatement: &Statement{
				Sql: "SELECT COUNT(`id`) FROM `test_model`;",
			},
		}, {
			name:     "with alias aggregate",
			selector: NewSelector[testModel](db).Select(Max("Age").As("max_age")),
			wantStatement: &Statement{
				Sql: "SELECT MAX(`age`) AS max_age FROM `test_model`;",
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

func TestSelector_FindOne(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, DBWithDialect(MySQLDialect), DBWithValueResolver(value.NewUnsafeResolver))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *Selector[testModel]
		wantRes  *testModel
		wantErr  error
	}{
		{
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
				rows.AddRow(1, "foo", 18, "bar")

				mock.ExpectQuery("SELECT *.").WillReturnRows(rows)
			},
			selector: NewSelector[testModel](db),
			wantRes: &testModel{
				Id:       1,
				Name:     "foo",
				Age:      18,
				NickName: &sql.NullString{String: "bar", Valid: true},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			var res *testModel
			res, err = tc.selector.FindOne(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}
