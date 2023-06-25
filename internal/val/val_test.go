package val

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jrmarcco/easy-orm/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type valArgs struct {
	ID       uint64
	Age      int
	Name     string
	NickName *sql.NullString
}

func testValWriteCols(t *testing.T, creator Creator) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	tcs := []struct {
		name     string
		entity   any
		mockRows func() *sqlmock.Rows
		wantRes  any
	}{
		{
			name:   "basic",
			entity: &valArgs{},
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "age", "name", "nick_name"})
				rows.AddRow(1, 18, "jrmarcco", "foo bar")
				return rows
			},
			wantRes: &valArgs{
				ID:       1,
				Age:      18,
				Name:     "jrmarcco",
				NickName: &sql.NullString{Valid: true, String: "foo bar"},
			},
		}, {
			name:   "out-of-order field",
			entity: &valArgs{},
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"age", "nick_name", "name", "id"})
				rows.AddRow(18, "foo bar", "jrmarcco", 1)
				return rows
			},
			wantRes: &valArgs{
				ID:       1,
				Age:      18,
				Name:     "jrmarcco",
				NickName: &sql.NullString{Valid: true, String: "foo bar"},
			},
		}, {
			name:   "partial field",
			entity: &valArgs{},
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "name"})
				rows.AddRow(1, "jrmarcco")
				return rows
			},
			wantRes: &valArgs{
				ID:   1,
				Name: "jrmarcco",
			},
		},
	}

	r := model.NewRegistry()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mockRows := tc.mockRows()
			mock.ExpectQuery("SELECT *.").WillReturnRows(mockRows)
			rows, err := mockDB.Query("SELECT *.")
			require.NoError(t, err)

			rows.Next()

			m, err := r.Get(tc.entity)
			require.NoError(t, err)

			val := creator(m, tc.entity)
			err = val.WriteCols(rows)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, tc.entity)
		})
	}
}

// go test -bench=BenchmarkVal_WriteCols -benchmem -benchtime=10000x
// goos: linux
// goarch: amd64
// pkg: github.com/jrmarcco/easy-orm/internal/val
// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkVal_WriteCols/reflect-16                  10000               972.3 ns/op           328 B/op         13 allocs/op
// BenchmarkVal_WriteCols/unsafe-16                   10000               541.5 ns/op           152 B/op          4 allocs/op
// PASS
// ok      github.com/jrmarcco/easy-orm/internal/val       0.025s
func BenchmarkVal_WriteCols(b *testing.B) {
	fn := func(b *testing.B, creator Creator) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(b, err)
		defer func(mockDB *sql.DB) {
			_ = mockDB.Close()
		}(mockDB)

		mockRows := mock.NewRows([]string{"id", "age", "name", "nick_name"})
		row := []driver.Value{1, 18, "jrmarcco", "foo bar"}

		for i := 0; i < b.N; i++ {
			mockRows.AddRow(row...)
		}

		mock.ExpectQuery("SELECT *.").WillReturnRows(mockRows)

		rows, err := mockDB.Query("SELECT *")

		r := model.NewRegistry()
		m, err := r.Get(&valArgs{})
		require.NoError(b, err)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			rows.Next()
			val := creator(m, &valArgs{})
			_ = val.WriteCols(rows)
		}
	}

	b.Run("reflect", func(b *testing.B) {
		fn(b, NewRefValWriter)
	})

	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeValWriter)
	})
}
