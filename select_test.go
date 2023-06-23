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
			name:    "basic * select with from db name",
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
			builder: NewSelector[selectorBuildArg](db).Select("Id", "FirstName"),
			wantStat: &Statement{
				SQL: "SELECT `id`,`first_name` FROM `selector_build_arg`;",
			},
		}, {
			name:    "assign invalid field select",
			builder: NewSelector[selectorBuildArg](db).Select("Id", "Invalid"),
			wantErr: errs.InvalidColumnFdErr("Invalid"),
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
