package orm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRawStat_Get(t *testing.T) {
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
		selector *RawStat[selectorBuildArg]
		wantRes  *selectorBuildArg
		wantErr  error
	}{
		{
			name: "error return",
			mockFunc: func() {
				mock.ExpectQuery("SELECT .*").
					WillReturnError(errors.New("this is an error msg"))
			},
			selector: NewRawStat[selectorBuildArg](
				db,
				`SELECT * FROM "selector_build_arg"`,
			),
			wantErr: errors.New("this is an error msg"),
		}, {
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: NewRawStat[selectorBuildArg](
				db,
				`SELECT * FROM "selector_build_arg"`,
			),
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

func TestRawStat_GetMulti(t *testing.T) {
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
		selector *RawStat[selectorBuildArg]
		wantRes  []*selectorBuildArg
		wantErr  error
	}{
		{
			name: "error return",
			mockFunc: func() {
				mock.ExpectQuery("SELECT .*").
					WillReturnError(errors.New("this is an error msg"))
			},
			selector: NewRawStat[selectorBuildArg](
				db,
				`SELECT * FROM "selector_build_arg"`,
			),
			wantErr: errors.New("this is an error msg"),
		}, {
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "age", "first_name", "last_name"})
				rows.AddRow(1, 18, "Tom", "Cat")
				rows.AddRow(2, 36, "Foo", "Bar")
				mock.ExpectQuery("SELECT .*").WillReturnRows(rows)
			},
			selector: NewRawStat[selectorBuildArg](
				db,
				`SELECT * FROM "selector_build_arg"`,
			),
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
			selector: NewRawStat[selectorBuildArg](
				db,
				`SELECT * FROM "selector_build_arg"`,
			),
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
