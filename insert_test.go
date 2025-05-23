package easyorm

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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
	db, err := OpenDB(&sql.DB{}, MySQLDialect)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		inserter *Inserter[insertTestModel]
		wantRes  *Statement
		wantErr  error
	}{
		{
			name: "basic",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id:      1,
				Age:     18,
				Name:    "foo",
				Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
				Balance: 100,
			}),
			wantRes: &Statement{
				SQL: "INSERT INTO `insert_test_model` (`id`, `age`, `name`, `email`, `balance`) VALUES (?, ?, ?, ?, ?);",
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL>", Valid: true}, float64(100),
				},
			},
			wantErr: nil,
		}, {
			name:     "without rows",
			inserter: NewInserter[insertTestModel](db),
			wantErr:  errs.ErrInsertWithoutRows,
		}, {
			name:     "with empty rows",
			inserter: NewInserter[insertTestModel](db).Rows(),
			wantErr:  errs.ErrInsertWithoutRows,
		}, {
			name: "with multiple rows",
			inserter: NewInserter[insertTestModel](db).Rows(
				&insertTestModel{
					Id:      1,
					Age:     18,
					Name:    "foo",
					Email:   &sql.NullString{String: "<EMAIL_1>", Valid: true},
					Balance: 100,
				},
				&insertTestModel{
					Id:      2,
					Age:     19,
					Name:    "bar",
					Email:   &sql.NullString{String: "<EMAIL_2>", Valid: true},
					Balance: 200,
				},
			),
			wantRes: &Statement{
				SQL: "INSERT INTO `insert_test_model` (`id`, `age`, `name`, `email`, `balance`) VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?);",
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL_1>", Valid: true}, float64(100),
					uint64(2), int8(19), "bar", &sql.NullString{String: "<EMAIL_2>", Valid: true}, float64(200),
				},
			},
		}, {
			name: "with columns",
			inserter: NewInserter[insertTestModel](db).Fields("Id", "Email").Rows(&insertTestModel{
				Id: 1,
				Email: &sql.NullString{
					String: "<EMAIL_1>",
					Valid:  true,
				},
			}, &insertTestModel{
				Id: 2,
				Email: &sql.NullString{
					String: "<EMAIL_2>",
					Valid:  true,
				},
			}),
			wantRes: &Statement{
				SQL: "INSERT INTO `insert_test_model` (`id`, `email`) VALUES (?, ?), (?, ?);",
				Args: []any{
					uint64(1), &sql.NullString{String: "<EMAIL_1>", Valid: true},
					uint64(2), &sql.NullString{String: "<EMAIL_2>", Valid: true},
				},
			},
		}, {
			name: "with on conflict",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id:   1,
				Age:  18,
				Name: "foo",
				Email: &sql.NullString{
					String: "<EMAIL>",
					Valid:  true,
				},
				Balance: 100,
			}).OnConflict().Update(Assign("Age", 19), Assign("Balance", 200)),
			wantRes: &Statement{
				SQL: "INSERT INTO `insert_test_model` (`id`, `age`, `name`, `email`, `balance`) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `age` = ?, `balance` = ?;",
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL>", Valid: true}, float64(100), 19, 200,
				},
			},
		}, {
			name: "with on conflict and assign column",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id:      1,
				Age:     18,
				Name:    "foo",
				Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
				Balance: 100,
			}).OnConflict().Update(Col("Age")),
			wantRes: &Statement{
				SQL: "INSERT INTO `insert_test_model` (`id`, `age`, `name`, `email`, `balance`) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `age` = VALUES(`age`);",
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL>", Valid: true}, float64(100),
				},
			},
		}, {
			name: "with invalid on conflict",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id: uint64(1),
			}).OnConflict().Update(Assign("Invalid", 19)),
			wantErr: errs.ErrInvalidField("Invalid"),
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

func TestInserter_OnConflict_Postgres(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, PostgresDialect)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		inserter *Inserter[insertTestModel]
		wantRes  *Statement
		wantErr  error
	}{
		{
			name: "without update",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id:      1,
				Age:     18,
				Name:    "foo",
				Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
				Balance: 100,
			}).OnConflict("Id").Update(),
			wantRes: &Statement{
				SQL: `INSERT INTO "insert_test_model" ("id", "age", "name", "email", "balance") VALUES ($1, $2, $3, $4, $5) ON CONFLICT ("id") DO NOTHING;`,
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL>", Valid: true}, float64(100),
				},
			},
		}, {
			name: "with update",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id:      1,
				Age:     18,
				Name:    "foo",
				Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
				Balance: 100,
			}).OnConflict("Id").Update(Assign("Age", 19), Assign("Balance", 200)),
			wantRes: &Statement{
				SQL: `INSERT INTO "insert_test_model" ("id", "age", "name", "email", "balance") VALUES ($1, $2, $3, $4, $5) ON CONFLICT ("id") DO UPDATE SET "age" = $6, "balance" = $7;`,
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL>", Valid: true}, float64(100), 19, 200,
				},
			},
		}, {
			name: "with update and assign column",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id:      1,
				Age:     18,
				Name:    "foo",
				Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
				Balance: 100,
			}).OnConflict("Id").Update(Col("Age"), Col("Balance")),
			wantRes: &Statement{
				SQL: `INSERT INTO "insert_test_model" ("id", "age", "name", "email", "balance") VALUES ($1, $2, $3, $4, $5) ON CONFLICT ("id") DO UPDATE SET "age" = EXCLUDED."age", "balance" = EXCLUDED."balance";`,
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL>", Valid: true}, float64(100),
				},
			},
		}, {
			name: "with multiple conflicts",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id:      1,
				Age:     18,
				Name:    "foo",
				Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
				Balance: 100,
			}).OnConflict("Id", "Name").Update(Assign("Age", 19), Assign("Balance", 200)),
			wantRes: &Statement{
				SQL: `INSERT INTO "insert_test_model" ("id", "age", "name", "email", "balance") VALUES ($1, $2, $3, $4, $5) ON CONFLICT ("id", "name") DO UPDATE SET "age" = $6, "balance" = $7;`,
				Args: []any{
					uint64(1), int8(18), "foo", &sql.NullString{String: "<EMAIL>", Valid: true}, float64(100), 19, 200,
				},
			},
		}, {
			name: "with invalid on conflict",
			inserter: NewInserter[insertTestModel](db).Rows(&insertTestModel{
				Id: uint64(1),
			}).OnConflict("Invalid").Update(Assign("Age", 19), Assign("Balance", 200)),
			wantErr: errs.ErrInvalidField("Invalid"),
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

func TestInserter_Exec(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, MySQLDialect)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		inserter *Inserter[insertTestModel]
		wantRes  int64
		wantErr  error
	}{
		{
			name: "basic",
			inserter: func() *Inserter[insertTestModel] {
				mock.ExpectExec("INSERT INTO `insert_test_model`.*").
					WillReturnResult(sqlmock.NewResult(1, 1))

				return NewInserter[insertTestModel](db).Rows(&insertTestModel{
					Id:      1,
					Age:     18,
					Name:    "foo",
					Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
					Balance: 100,
				})
			}(),
			wantRes: 1,
		}, {
			name: "db error",
			inserter: func() *Inserter[insertTestModel] {
				mock.ExpectExec("INSERT INTO `insert_test_model`.*").
					WillReturnError(errors.New("db error"))

				return NewInserter[insertTestModel](db).Rows(&insertTestModel{
					Id:      1,
					Age:     18,
					Name:    "foo",
					Email:   &sql.NullString{String: "<EMAIL>", Valid: true},
					Balance: 100,
				})
			}(),
			wantErr: errors.New("db error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.inserter.Exec(context.Background())
			assert.Equal(t, tc.wantErr, res.Err())

			if res.Err() == nil {
				assert.Equal(t, tc.wantRes, res.RowsAffected())
			}
		})
	}
}
