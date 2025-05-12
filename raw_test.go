package easyorm

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type rawTestModel struct {
	Id   uint64
	Name *sql.NullString
}

func TestRaw_FindOne(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, DBWithDialect(MySQLDialect))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		raw      *Raw[rawTestModel]
		mockFunc func()
		wantRes  *rawTestModel
		wantErr  error
	}{
		{
			name: "basic",
			raw: NewRaw[rawTestModel](
				db,
				`SELECT * FROM raw_test_model`,
			),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"})
				rows.AddRow(1, "foo")

				mock.ExpectQuery("SELECT *.").WillReturnRows(rows)
			},
			wantRes: &rawTestModel{
				Id: 1,
				Name: &sql.NullString{
					String: "foo",
					Valid:  true,
				},
			},
		}, {
			name: "with args",
			raw: NewRaw[rawTestModel](
				db,
				`SELECT * FROM raw_test_model WHERE id = ? AND name = ?`,
				1, "foo",
			),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"})
				rows.AddRow(1, "foo")

				mock.ExpectQuery("SELECT *.").WithArgs(1, "foo").WillReturnRows(rows)
			},
			wantRes: &rawTestModel{
				Id: 1,
				Name: &sql.NullString{
					String: "foo",
					Valid:  true,
				},
			},
		}, {
			name: "returns error",
			raw: NewRaw[rawTestModel](
				db,
				`SELECT * FROM raw_test_model`,
			),
			mockFunc: func() {

				mock.ExpectQuery("SELECT *.").WillReturnError(errors.New("mock error"))
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			var res *rawTestModel
			res, err = tc.raw.FindOne(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}

func TestRaw_FindMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, DBWithDialect(MySQLDialect))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		raw      *Raw[rawTestModel]
		mockFunc func()
		wantRes  []*rawTestModel
		wantErr  error
	}{
		{
			name: "basic",
			raw: NewRaw[rawTestModel](
				db,
				`SELECT * FROM raw_test_model`,
			),
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name"})
				rows.AddRow(1, "foo")
				rows.AddRow(2, "bar")

				mock.ExpectQuery("SELECT *.").WillReturnRows(rows)
			},
			wantRes: []*rawTestModel{
				{
					Id: 1,
					Name: &sql.NullString{
						String: "foo",
						Valid:  true,
					},
				}, {
					Id: 2,
					Name: &sql.NullString{
						String: "bar",
						Valid:  true,
					},
				},
			},
		}, {
			name: "returns error",
			raw: NewRaw[rawTestModel](
				db,
				`SELECT * FROM raw_test_model`,
			),
			mockFunc: func() {
				mock.ExpectQuery("SELECT *.").WillReturnError(errors.New("mock error"))
			},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			var res []*rawTestModel
			res, err = tc.raw.FindMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}

func TestRaw_Exec(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, DBWithDialect(MySQLDialect))
	require.NoError(t, err)

	tcs := []struct {
		name    string
		raw     *Raw[rawTestModel]
		wantRes int64
		wantErr error
	}{
		{
			name: "basic",
			raw: func() *Raw[rawTestModel] {
				mock.ExpectExec("INSERT INTO raw_test_model.*").
					WillReturnResult(sqlmock.NewResult(1, 1))

				return NewRaw[rawTestModel](
					db,
					"INSERT INTO raw_test_model (id, name) VALUES (?, ?)",
					uint64(1),
					&sql.NullString{
						String: "foo",
						Valid:  true,
					},
				)
			}(),
			wantRes: 1,
		}, {
			name: "db error",
			raw: func() *Raw[rawTestModel] {
				mock.ExpectExec("INSERT INTO raw_test_model.*").
					WillReturnError(errors.New("mock error"))

				return NewRaw[rawTestModel](
					db,
					"INSERT INTO raw_test_model (id, name) VALUES (?, ?)",
					uint64(1),
					&sql.NullString{
						String: "foo",
						Valid:  true,
					})
			}(),
			wantErr: errors.New("mock error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.raw.Exec(context.Background())
			assert.Equal(t, tc.wantErr, res.Err())

			if res.Err() == nil {
				assert.Equal(t, tc.wantRes, res.RowsAffected())
			}
		})
	}

}
