package slog

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	easyorm "github.com/JrMarcco/easy-orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type slogTestModel struct {
	Id       uint64
	Name     string
	Age      int8
	NickName *sql.NullString
}

func TestMiddlewareBuilder_Build(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	opt := easyorm.DBWithMiddlewareChain(easyorm.MiddlewareChain{
		NewMiddlewareBuilder().Build(),
	})

	db, err := easyorm.OpenDB(mockDB, easyorm.PostgresDialect, opt)
	require.NoError(t, err)

	selector := easyorm.NewSelector[slogTestModel](db).Where(easyorm.Col("Id").Eq(1))

	mockRows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
	mockRows.AddRow(1, "foo", 18, "bar")
	mock.ExpectQuery("SELECT .*").WillReturnRows(mockRows)

	one, err := selector.FindOne(context.Background())
	require.NoError(t, err)
	assert.Equal(t, one, &slogTestModel{
		Id:       1,
		Name:     "foo",
		Age:      18,
		NickName: &sql.NullString{String: "bar", Valid: true},
	})

	inserter := easyorm.NewInserter[slogTestModel](db)
	mock.ExpectExec("INSERT INTO .*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	res := inserter.Rows(&slogTestModel{
		Id:       1,
		Name:     "foo",
		Age:      18,
		NickName: &sql.NullString{String: "bar", Valid: true},
	}).Exec(context.Background())
	assert.NoError(t, res.Err())
	assert.Equal(t, int64(1), res.RowsAffected())

	deleter := easyorm.NewDeleter[slogTestModel](db)
	mock.ExpectExec("DELETE FROM .*").
		WillReturnResult(sqlmock.NewResult(1, 1))
	res = deleter.Where(easyorm.Col("Id").Eq(1)).Exec(context.Background())
	assert.NoError(t, res.Err())
	assert.Equal(t, int64(1), res.RowsAffected())

}
