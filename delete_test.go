package easyorm

import (
	"database/sql"
	"testing"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type deleteTestModel struct {
	Id    uint64
	Age   int8
	Name  string
	Email *sql.NullString
}

func TestDeleter_Build(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, MySQLDialect)
	require.NoError(t, err)

	tcs := []struct {
		name    string
		deleter *Deleter[deleteTestModel]
		wantRes *Statement
		wantErr error
	}{
		{
			name:    "basic",
			deleter: NewDeleter[deleteTestModel](db),
			wantRes: &Statement{
				SQL: "DELETE FROM `delete_test_model`;",
			},
			wantErr: nil,
		}, {
			name:    "with where",
			deleter: NewDeleter[deleteTestModel](db).Where(Col("Id").Eq(1)),
			wantRes: &Statement{
				SQL:  "DELETE FROM `delete_test_model` WHERE `id` = ?;",
				Args: []any{1},
			},
			wantErr: nil,
		}, {
			name:    "with empty where",
			deleter: NewDeleter[deleteTestModel](db).Where(),
			wantRes: &Statement{
				SQL: "DELETE FROM `delete_test_model`;",
			},
			wantErr: nil,
		}, {
			name:    "with where not",
			deleter: NewDeleter[deleteTestModel](db).Where(Col("Id").Eq(1).Not()),
			wantRes: &Statement{
				SQL:  "DELETE FROM `delete_test_model` WHERE NOT (`id` = ?);",
				Args: []any{1},
			},
		}, {
			name:    "with invalid where",
			deleter: NewDeleter[deleteTestModel](db).Where(Col("ID").Eq(1)),
			wantErr: errs.ErrInvalidField("ID"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			statement, err := tc.deleter.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, statement)
			}
		})
	}
}
