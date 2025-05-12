package value

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/JrMarcco/easy-orm/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type vrTestModel struct {
	Id       uint64
	Age      int8
	Name     string
	NickName *sql.NullString
}

func writeColumnsTestFunc(t *testing.T, rc ResolverCreator) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	r := model.NewRegistry()

	tcs := []struct {
		name     string
		entity   any
		mockRows *sqlmock.Rows
		wantRes  any
		wantErr  error
	}{
		{
			name:   "basic",
			entity: &vrTestModel{},
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "age", "name", "nick_name"})
				rows.AddRow(1, 18, "foo", "bar")
				return rows
			}(),
			wantRes: &vrTestModel{
				Id:       1,
				Age:      18,
				Name:     "foo",
				NickName: &sql.NullString{String: "bar", Valid: true},
			},
		}, {
			name:   "partial field",
			entity: &vrTestModel{},
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"id", "nick_name"})
				rows.AddRow(1, "bar")
				return rows
			}(),
			wantRes: &vrTestModel{
				Id: 1,
				NickName: &sql.NullString{
					String: "bar",
					Valid:  true,
				},
			},
		}, {
			name:   "out-of-order field",
			entity: &vrTestModel{},
			mockRows: func() *sqlmock.Rows {
				rows := sqlmock.NewRows([]string{"nick_name", "id", "name", "age"})
				rows.AddRow("bar", 1, "foo", 18)
				return rows
			}(),
			wantRes: &vrTestModel{
				Id:   1,
				Name: "foo",
				Age:  18,
				NickName: &sql.NullString{
					String: "bar",
					Valid:  true,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectQuery("SELECT *.").WillReturnRows(tc.mockRows)
			rows, err := mockDB.Query("SELECT *.")
			require.NoError(t, err)

			m, err := r.GetModel(tc.entity)
			require.NoError(t, err)

			resolver := rc(m, tc.entity)

			_ = rows.Next()
			err = resolver.WriteColumns(rows)
			require.NoError(t, err)

			assert.Equal(t, tc.wantRes, tc.entity)
		})
	}

}

type vrBenchMarkModel struct {
	Id       uint64
	Age      int8
	Name     string
	NickName *sql.NullString
	IDCardNo string
	Address  string
	Phone    string
	Email    string
}

// go test -bench=BenchmarkWriteColumns -benchme
// goos: linux
// goarch: amd64
// pkg: github.com/JrMarcco/easy-orm/internal/value
// cpu: 13th Gen Intel(R) core(TM) i7-13700KF
// BenchmarkWriteColumns/reflect-24                 1347993               897.6 ns/op           592 B/op         13 allocs/op
// BenchmarkWriteColumns/unsafe-24                  3359084               358.5 ns/op           280 B/op          4 allocs/op
func BenchmarkWriteColumns(b *testing.B) {
	fn := func(b *testing.B, rc ResolverCreator) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(b, err)
		defer func() {
			_ = mockDB.Close()
		}()

		mockRows := sqlmock.NewRows([]string{"id", "age", "name", "nick_name", "id_card_no", "address", "phone", "email"})
		row := []driver.Value{1, 18, "foo", "bar", "1234567890", "address", "1234567890", "<EMAIL>"}

		for i := 0; i < b.N; i++ {
			mockRows.AddRow(row...)
		}

		mock.ExpectQuery("SELECT *.").WillReturnRows(mockRows)

		sqlRows, err := mockDB.Query("SELECT *.")
		require.NoError(b, err)

		r := model.NewRegistry()
		m, err := r.GetModel(&vrBenchMarkModel{})
		require.NoError(b, err)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sqlRows.Next()
			resolver := rc(m, &vrBenchMarkModel{})
			_ = resolver.WriteColumns(sqlRows)
		}
	}

	b.Run("reflect", func(b *testing.B) {
		fn(b, NewReflectResolver)
	})

	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeResolver)
	})
}
