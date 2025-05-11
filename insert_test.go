package easyorm

import (
	"database/sql"
	"testing"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type insertTestModel struct {
	Id      uint64
	Age     int8
	Name    string
	Email   *sql.NullString
	Balance float64
}

func TestInserter_Build(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, DBWithDialect(MySQLDialect))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		inserter *Inserter[insertTestModel]
		wantRes  *Statement
		wantErr  error
	}{
		{
			name: "basic",
			inserter: NewInserter[insertTestModel](db).Insert(&insertTestModel{
				Id:      1,
				Age:     18,
				Name:    "foo",
				Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
				Balance: 100,
			}),
			wantRes: &Statement{
				SQL: "INSERT INTO `insert_test_model` (`id`, `age`, `name`, `email`, `balance`) VALUES (?, ?, ?, ?, ?);",
			},
			wantErr: nil,
		}, {
			name:     "without rows",
			inserter: NewInserter[insertTestModel](db),
			wantErr:  errs.ErrInsertWithoutRows,
		}, {
			name:     "with empty rows",
			inserter: NewInserter[insertTestModel](db).Insert(),
			wantErr:  errs.ErrInsertWithoutRows,
		}, {
			name: "with multiple rows",
			inserter: NewInserter[insertTestModel](db).Insert(
				&insertTestModel{
					Id:      1,
					Age:     18,
					Name:    "foo",
					Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
					Balance: 100,
				},
				&insertTestModel{
					Id:      2,
					Age:     19,
					Name:    "bar",
					Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
					Balance: 200,
				},
			),
			wantRes: &Statement{
				SQL: "INSERT INTO `insert_test_model` (`id`, `age`, `name`, `email`, `balance`) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?);",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			statement, err := tc.inserter.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, statement)
			}
		})
	}
}
