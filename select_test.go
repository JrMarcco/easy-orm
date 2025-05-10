package easyorm

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type selectTestModel struct {
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
			selector: NewSelector[selectTestModel](db),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model`;",
			},
			wantErr: nil,
		}, {
			name:     "with where",
			selector: NewSelector[selectTestModel](db).Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "with where not",
			selector: NewSelector[selectTestModel](db).Where(Col("Id").Eq(1).Not()),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE NOT (`id` = ?);",
				Args: []any{
					1,
				},
			},
		}, {
			name: "with where and",
			selector: NewSelector[selectTestModel](db).Where(
				Col("Id").Eq(1),
				Col("Name").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE (`id` = ?) AND (`name` = ?);",
				Args: []any{
					1, "jrmarcco",
				},
			},
		}, {
			name: "with where or",
			selector: NewSelector[selectTestModel](db).Where(
				Col("Age").Ge(18).
					Or(
						Col("Age").Lt(12),
					),
				Col("NickName").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE ((`age` >= ?) OR (`age` < ?)) AND (`nick_name` = ?);",
				Args: []any{
					18, 12, "jrmarcco",
				},
			},
		}, {
			name: "with invalid column",
			selector: NewSelector[selectTestModel](db).Where(
				Col("InvalidColumn").Eq(1),
			),
			wantErr: errs.ErrInvalidField("InvalidColumn"),
		}, {
			name: "with basic selectable",
			selector: NewSelector[selectTestModel](db).
				Select(
					Col("Id"),
					Col("Name"),
					Col("Age"),
				).
				Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				SQL: "SELECT `id`, `name`, `age` FROM `select_test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name: "with alias selectable",
			selector: NewSelector[selectTestModel](db).
				Select(
					Col("Id"),
					Col("Name").As("user_name"),
					Col("Age"),
				).
				Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				SQL: "SELECT `id`, `name` AS `user_name`, `age` FROM `select_test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "with count aggregate",
			selector: NewSelector[selectTestModel](db).Select(Count("Id")),
			wantStatement: &Statement{
				SQL: "SELECT COUNT(`id`) FROM `select_test_model`;",
			},
		}, {
			name:     "with sum aggregate",
			selector: NewSelector[selectTestModel](db).Select(Sum("Age")),
			wantStatement: &Statement{
				SQL: "SELECT SUM(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with max aggregate",
			selector: NewSelector[selectTestModel](db).Select(Max("Age")),
			wantStatement: &Statement{
				SQL: "SELECT MAX(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with min aggregate",
			selector: NewSelector[selectTestModel](db).Select(Min("Age")),
			wantStatement: &Statement{
				SQL: "SELECT MIN(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with avg aggregate",
			selector: NewSelector[selectTestModel](db).Select(Avg("Age")),
			wantStatement: &Statement{
				SQL: "SELECT AVG(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with alias aggregate",
			selector: NewSelector[selectTestModel](db).Select(Max("Age").As("max_age")),
			wantStatement: &Statement{
				SQL: "SELECT MAX(`age`) AS `max_age` FROM `select_test_model`;",
			},
		}, {
			name: "with invalid aggregate",
			selector: NewSelector[selectTestModel](db).Select(
				Max("InvalidColumn"),
			),
			wantErr: errs.ErrInvalidField("InvalidColumn"),
		}, {
			name: "with multiple aggregates",
			selector: NewSelector[selectTestModel](db).Select(
				Count("Id"),
				Sum("Age"),
				Max("Age"),
				Min("Age"),
				Avg("Age"),
			),
			wantStatement: &Statement{
				SQL: "SELECT COUNT(`id`), SUM(`age`), MAX(`age`), MIN(`age`), AVG(`age`) FROM `select_test_model`;",
			},
		}, {
			name: "with raw expression in where",
			selector: NewSelector[selectTestModel](db).
				Where(
					RawExpr("`id` = (`age` + ?)", 1),
				),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` = (`age` + ?);",
				Args: []any{
					1,
				},
			},
			wantErr: nil,
		}, {
			name: "with raw expression in predicate",
			selector: NewSelector[selectTestModel](db).
				Where(
					Col("Id").Eq(RawExpr("`id` + ?", 1)),
				),
			wantStatement: &Statement{
				SQL:  "SELECT * FROM `select_test_model` WHERE `id` = (`id` + ?);",
				Args: []any{1},
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

	db, err := OpenDB(mockDB, DBWithDialect(MySQLDialect))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *Selector[selectTestModel]
		wantRes  *selectTestModel
		wantErr  error
	}{
		{
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
				rows.AddRow(1, "foo", 18, "bar")

				mock.ExpectQuery("SELECT *.").WillReturnRows(rows)
			},
			selector: NewSelector[selectTestModel](db),
			wantRes: &selectTestModel{
				Id:       1,
				Name:     "foo",
				Age:      18,
				NickName: &sql.NullString{String: "bar", Valid: true},
			},
		}, {
			name: "with args",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
				rows.AddRow(1, "foo", 18, "bar")

				mock.ExpectQuery("SELECT *.").WithArgs(1, "foo").WillReturnRows(rows)
			},
			selector: NewSelector[selectTestModel](db).
				Where(
					Col("Id").Eq(1),
					Col("Name").Eq("foo"),
				),
			wantRes: &selectTestModel{
				Id:       1,
				Name:     "foo",
				Age:      18,
				NickName: &sql.NullString{String: "bar", Valid: true},
			},
		}, {
			name: "returns error",
			mockFunc: func() {
				mock.ExpectQuery("SELECT *.").WillReturnError(errors.New("mock error"))
			},
			selector: NewSelector[selectTestModel](db),
			wantErr:  errors.New("mock error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			var res *selectTestModel
			res, err = tc.selector.FindOne(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}

func TestSelector_FindMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, DBWithDialect(MySQLDialect))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *Selector[selectTestModel]
		wantRes  []*selectTestModel
		wantErr  error
	}{
		{
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
				rows.AddRow(1, "foo", 18, "bar")
				rows.AddRow(2, "bar", 12, "foo")
				rows.AddRow(3, "baz", 10, "baz")

				mock.ExpectQuery("SELECT *.").WillReturnRows(rows)
			},
			selector: NewSelector[selectTestModel](db),
			wantRes: []*selectTestModel{
				{
					Id:       1,
					Name:     "foo",
					Age:      18,
					NickName: &sql.NullString{String: "bar", Valid: true},
				}, {
					Id:       2,
					Name:     "bar",
					Age:      12,
					NickName: &sql.NullString{String: "foo", Valid: true},
				}, {
					Id:       3,
					Name:     "baz",
					Age:      10,
					NickName: &sql.NullString{String: "baz", Valid: true},
				},
			},
		}, {
			name: "returns error",
			mockFunc: func() {
				mock.ExpectQuery("SELECT *.").WillReturnError(errors.New("mock error"))
			},
			selector: NewSelector[selectTestModel](db),
			wantErr:  errors.New("mock error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			var res []*selectTestModel
			res, err = tc.selector.FindMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}
