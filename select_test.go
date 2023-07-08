package orm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelector_Build(t *testing.T) {

	db, err := OpenDB(&sql.DB{})
	require.NoError(t, err)

	tcs := []struct {
		name     string
		builder  StatBuilder
		wantStat *Statement
		wantErr  error
	}{
		{
			name:    "basic * select without from",
			builder: NewSelector[selectorBuildArg](db),
			wantStat: &Statement{
				SQL: "SELECT * FROM `selector_build_arg`;",
			},
		},
		{
			name:    "basic * select with from",
			builder: NewSelector[selectorBuildArg](db).From("test_model"),
			wantStat: &Statement{
				SQL: "SELECT * FROM `test_model`;",
			},
		}, {
			name:    "basic * select with empty from",
			builder: NewSelector[selectorBuildArg](db).From(""),
			wantStat: &Statement{
				SQL: "SELECT * FROM `selector_build_arg`;",
			},
		}, {
			name:    "basic * select with from db fdName",
			builder: NewSelector[selectorBuildArg](db).From("test_db.test_model"),
			wantStat: &Statement{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		}, {
			name:    "empty where",
			builder: NewSelector[selectorBuildArg](db).Where(),
			wantStat: &Statement{
				SQL: "SELECT * FROM `selector_build_arg`;",
			},
		}, {
			name:    "single predicate where",
			builder: NewSelector[selectorBuildArg](db).Where(Col("Age").Eq(18)),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE `age` = ?;",
				Args: []any{18},
			},
		}, {
			name:    "not predicate where",
			builder: NewSelector[selectorBuildArg](db).Where(Not(Col("Age").Eq(18))),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE NOT (`age` = ?);",
				Args: []any{18},
			},
		}, {
			name: "not & and predicate where",
			builder: NewSelector[selectorBuildArg](db).Where(
				Not(
					Col("Age").Eq(18).And(Col("Id").Eq(1)),
				),
			),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE NOT ((`age` = ?) AND (`id` = ?));",
				Args: []any{18, 1},
			},
		}, {
			name: "not & or predicate where",
			builder: NewSelector[selectorBuildArg](db).Where(
				Not(
					Col("Id").Gt(100).Or(Col("Age").Lt(18)),
				),
			),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE NOT ((`id` > ?) OR (`age` < ?));",
				Args: []any{100, 18},
			},
		}, {
			name:    "invalid type",
			builder: NewSelector[selectorBuildArg](db).Where(Col("Invalid").Eq("test")),
			wantErr: errs.InvalidColumnFdErr("Invalid"),
		}, {
			name:    "assign field select",
			builder: NewSelector[selectorBuildArg](db).Select(Col("Id"), Col("FirstName")),
			wantStat: &Statement{
				SQL: "SELECT `id`,`first_name` FROM `selector_build_arg`;",
			},
		}, {
			name:    "assign invalid field select",
			builder: NewSelector[selectorBuildArg](db).Select(Col("Id"), Col("Invalid")),
			wantErr: errs.InvalidColumnFdErr("Invalid"),
		}, {
			name:    "avg aggregate func select",
			builder: NewSelector[selectorBuildArg](db).Select(Avg("Age")),
			wantStat: &Statement{
				SQL: "SELECT AVG(`age`) FROM `selector_build_arg`;",
			},
		}, {
			name:    "sum aggregate func select",
			builder: NewSelector[selectorBuildArg](db).Select(Sum("Age")),
			wantStat: &Statement{
				SQL: "SELECT SUM(`age`) FROM `selector_build_arg`;",
			},
		}, {
			name:    "count aggregate func select",
			builder: NewSelector[selectorBuildArg](db).Select(Count("Age")),
			wantStat: &Statement{
				SQL: "SELECT COUNT(`age`) FROM `selector_build_arg`;",
			},
		}, {
			name:    "max aggregate func select",
			builder: NewSelector[selectorBuildArg](db).Select(Max("FirstName")),
			wantStat: &Statement{
				SQL: "SELECT MAX(`first_name`) FROM `selector_build_arg`;",
			},
		}, {
			name:    "min aggregate func select",
			builder: NewSelector[selectorBuildArg](db).Select(Min("Id")),
			wantStat: &Statement{
				SQL: "SELECT MIN(`id`) FROM `selector_build_arg`;",
			},
		}, {
			name:    "multi aggregate func select",
			builder: NewSelector[selectorBuildArg](db).Select(Max("FirstName"), Min("Id")),
			wantStat: &Statement{
				SQL: "SELECT MAX(`first_name`),MIN(`id`) FROM `selector_build_arg`;",
			},
		}, {
			name:    "invalid field aggregate func select",
			builder: NewSelector[selectorBuildArg](db).Select(Avg("Invalid")),
			wantErr: errs.InvalidColumnFdErr("Invalid"),
		}, {
			name:    "raw expression",
			builder: NewSelector[selectorBuildArg](db).Select(Raw("COUNT(`Id`) AS id_count")),
			wantStat: &Statement{
				SQL: "SELECT COUNT(`Id`) AS id_count FROM `selector_build_arg`;",
			},
		}, {
			name: "raw expression in where",
			builder: NewSelector[selectorBuildArg](db).Where(
				Raw("`first_name` LIKE %?%", "jrmarcco").AsPredicate(),
			),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE `first_name` LIKE %?%;",
				Args: []any{"jrmarcco"},
			},
		}, {
			name: "row expression in predicate",
			builder: NewSelector[selectorBuildArg](db).Where(
				Col("Id").Eq(Raw("`age` + ?", 10000).AsPredicate()),
			),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE `id` = (`age` + ?);",
				Args: []any{10000},
			},
		}, {
			name: "alias",
			builder: NewSelector[selectorBuildArg](db).Select(
				Col("FirstName").As("aliasName"),
			),
			wantStat: &Statement{
				SQL: "SELECT `first_name` AS `aliasName` FROM `selector_build_arg`;",
			},
		}, {
			name: "multi alias",
			builder: NewSelector[selectorBuildArg](db).Select(
				Col("Id").As("aliasId"),
				Col("FirstName").As("aliasName"),
			),
			wantStat: &Statement{
				SQL: "SELECT `id` AS `aliasId`,`first_name` AS `aliasName` FROM `selector_build_arg`;",
			},
		}, {
			name: "alias in aggregate",
			builder: NewSelector[selectorBuildArg](db).Select(
				Avg("Age").As("avgAge"),
			),
			wantStat: &Statement{
				SQL: "SELECT AVG(`age`) AS `avgAge` FROM `selector_build_arg`;",
			},
		}, {
			name: "multi alias in aggregate",
			builder: NewSelector[selectorBuildArg](db).Select(
				Avg("Age").As("avgAge"),
				Sum("Age").As("sumAge"),
			),
			wantStat: &Statement{
				SQL: "SELECT AVG(`age`) AS `avgAge`,SUM(`age`) AS `sumAge` FROM `selector_build_arg`;",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stat, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantStat, stat)
			}
		})
	}

}

type selectorBuildArg struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  *sql.NullString
}

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *Selector[selectorBuildArg]
		wantRes  *selectorBuildArg
		wantErr  error
	}{
		{
			name:     "invalid query",
			selector: NewSelector[selectorBuildArg](db).Where(Col("Invalid").Eq("...")),
			wantErr:  errs.InvalidColumnFdErr("Invalid"),
		}, {
			name: "error return",
			mockFunc: func() {
				mock.ExpectQuery("SELECT .*").
					WillReturnError(errors.New("this is an error msg"))
			},
			selector: NewSelector[selectorBuildArg](db).Where(),
			wantErr:  errors.New("this is an error msg"),
		}, {
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: NewSelector[selectorBuildArg](db).Where(Col("Id").Eq(1)),
			wantRes: &selectorBuildArg{
				Id:        1,
				Age:       18,
				FirstName: "Tom",
				LastName: &sql.NullString{
					Valid:  true,
					String: "Cat",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			res, err := tc.selector.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}

func TestSelector_GetMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *Selector[selectorBuildArg]
		wantRes  []*selectorBuildArg
		wantErr  error
	}{
		{
			name:     "invalid query",
			selector: NewSelector[selectorBuildArg](db).Where(Col("Invalid").Eq("...")),
			wantErr:  errs.InvalidColumnFdErr("Invalid"),
		}, {
			name: "error return",
			mockFunc: func() {
				mock.ExpectQuery("SELECT .*").
					WillReturnError(errors.New("this is an error msg"))
			},
			selector: NewSelector[selectorBuildArg](db).Where(),
			wantErr:  errors.New("this is an error msg"),
		}, {
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				rows.AddRow(2, 36, "Foo", "Bar")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: NewSelector[selectorBuildArg](db).Where(Col("Id").Gt(0)),
			wantRes: []*selectorBuildArg{
				{
					Id:        1,
					Age:       18,
					FirstName: "Tom",
					LastName: &sql.NullString{
						Valid:  true,
						String: "Cat",
					},
				}, {
					Id:        2,
					Age:       36,
					FirstName: "Foo",
					LastName: &sql.NullString{
						Valid:  true,
						String: "Bar",
					},
				},
			},
		}, {
			name: "proportion columns",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				rows.AddRow(2, 36, "Foo", "Bar")
				rows.AddRow(3, 30, "", "jrmarcco")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: NewSelector[selectorBuildArg](db).Where(Col("Id").Gt(0)),
			wantRes: []*selectorBuildArg{
				{
					Id:        1,
					Age:       18,
					FirstName: "Tom",
					LastName: &sql.NullString{
						Valid:  true,
						String: "Cat",
					},
				}, {
					Id:        2,
					Age:       36,
					FirstName: "Foo",
					LastName: &sql.NullString{
						Valid:  true,
						String: "Bar",
					},
				}, {
					Id:  3,
					Age: 30,
					LastName: &sql.NullString{
						Valid:  true,
						String: "jrmarcco",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			res, err := tc.selector.GetMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}
