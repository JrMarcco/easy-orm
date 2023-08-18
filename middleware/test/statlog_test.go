package test

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	orm "github.com/jrmarcco/easy-orm"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/middleware/statlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

type selectorBuildArg struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  *sql.NullString
}

var globalStat string
var globalArgs []any

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	logMdlBuilder := statlog.NewBuilder(statlog.BuilderWithLogFunc(func(stat *orm.Statement) {
		log.Printf("statement: %s, args: %v \n", stat.SQL, stat.Args)
		globalStat = stat.SQL
		globalArgs = stat.Args
	}))

	db, err := orm.OpenDB(mockDB, orm.DBWithDialect(orm.MySqlDialect), orm.DBWithMdls(logMdlBuilder.Build()))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *orm.Selector[selectorBuildArg]
		wantStat string
		wantArgs []any
		wantErr  error
	}{
		{
			name:     "invalid query",
			selector: orm.NewSelector[selectorBuildArg](db).Where(orm.Col("Invalid").Eq("...")),
			wantErr:  errs.InvalidColumnFdErr("Invalid"),
		}, {
			name: "error return",
			mockFunc: func() {
				mock.ExpectQuery("SELECT .*").
					WillReturnError(errors.New("this is an error msg"))
			},
			selector: orm.NewSelector[selectorBuildArg](db),
			wantStat: "SELECT * FROM `selector_build_arg`;",
			wantErr:  errors.New("this is an error msg"),
		}, {
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: orm.NewSelector[selectorBuildArg](db).Where(orm.Col("Id").Eq(1)),
			wantStat: "SELECT * FROM `selector_build_arg` WHERE `id` = ?;",
			wantArgs: []any{1},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			globalStat = ""
			globalArgs = nil

			_, err := tc.selector.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)

			assert.Equal(t, tc.wantStat, globalStat)
			assert.Equal(t, tc.wantArgs, globalArgs)
		})
	}
}

func TestSelector_GetMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	logMdlBuilder := statlog.NewBuilder(statlog.BuilderWithLogFunc(func(stat *orm.Statement) {
		log.Printf("statement: %s, args: %v \n", stat.SQL, stat.Args)
		globalStat = stat.SQL
		globalArgs = stat.Args
	}))

	db, err := orm.OpenDB(mockDB, orm.DBWithDialect(orm.MySqlDialect), orm.DBWithMdls(logMdlBuilder.Build()))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *orm.Selector[selectorBuildArg]
		wantStat string
		wantArgs []any
		wantErr  error
	}{
		{
			name:     "invalid query",
			selector: orm.NewSelector[selectorBuildArg](db).Where(orm.Col("Invalid").Eq("...")),
			wantErr:  errs.InvalidColumnFdErr("Invalid"),
		}, {
			name: "error return",
			mockFunc: func() {
				mock.ExpectQuery("SELECT .*").
					WillReturnError(errors.New("this is an error msg"))
			},
			selector: orm.NewSelector[selectorBuildArg](db).Where(),
			wantStat: "SELECT * FROM `selector_build_arg`;",
			wantErr:  errors.New("this is an error msg"),
		}, {
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				rows.AddRow(2, 36, "Foo", "Bar")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: orm.NewSelector[selectorBuildArg](db).Where(orm.Col("Id").Gt(0)),
			wantStat: "SELECT * FROM `selector_build_arg` WHERE `id` > ?;",
			wantArgs: []any{0},
		}, {
			name: "proportion columns",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				rows.AddRow(2, 36, "Foo", "Bar")
				rows.AddRow(3, 30, "", "jrmarcco")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: orm.NewSelector[selectorBuildArg](db).Where(orm.Col("Id").Gt(0)),
			wantStat: "SELECT * FROM `selector_build_arg` WHERE `id` > ?;",
			wantArgs: []any{0},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			if tc.mockFunc != nil {
				tc.mockFunc()
			}

			globalStat = ""
			globalArgs = nil

			_, err := tc.selector.GetMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)

			assert.Equal(t, tc.wantStat, globalStat)
			assert.Equal(t, tc.wantArgs, globalArgs)
		})
	}
}
