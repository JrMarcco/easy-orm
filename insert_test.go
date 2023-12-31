package orm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type inserterBuildArg struct {
	Id       uint64
	Name     string
	NickName *sql.NullString
	Balance  int64
}

type invalidAssign struct {
}

func (i invalidAssign) assign() {}
func TestInserter_Build(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, DBWithDialect(MySqlDialect))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		inserter StatBuilder
		wantStat *Statement
		wantErr  error
	}{
		{
			name:     "empty row",
			inserter: NewInserter[inserterBuildArg](db),
			wantErr:  errs.ErrEmptyInsertRow,
		}, {
			name: "single row",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`,`nick_name`,`balance`) VALUES (?,?,?,?);",
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100),
				},
			},
		}, {
			name: "multi row",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
				&inserterBuildArg{
					Id:   uint64(2),
					Name: "tom cat",
					NickName: &sql.NullString{
						Valid:  true,
						String: "cat",
					},
					Balance: int64(1),
				},
			),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`,`nick_name`,`balance`) VALUES (?,?,?,?),(?,?,?,?);",
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100),
					uint64(2), "tom cat", &sql.NullString{Valid: true, String: "cat"}, int64(1),
				},
			},
		}, {
			name: "specify insert column field",
			inserter: NewInserter[inserterBuildArg](db).ColFd("Id", "Name").Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
				},
			),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`) VALUES (?,?);",
				Args: []any{
					uint64(1), "jrmarcco",
				},
			},
		}, {
			name: "specify invalid insert column field",
			inserter: NewInserter[inserterBuildArg](db).ColFd("Id", "Invalid").Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
				},
			),
			wantErr: errs.ErrInvalidColumnFd("Invalid"),
		}, {
			name: "specify insert column field with multi row",
			inserter: NewInserter[inserterBuildArg](db).ColFd("Id", "Name").Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
				},
				&inserterBuildArg{
					Id:   uint64(2),
					Name: "tom cat",
				},
			),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`) VALUES (?,?),(?,?);",
				Args: []any{
					uint64(1), "jrmarcco",
					uint64(2), "tom cat",
				},
			},
		}, {
			name: "update column on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts().Update(Col("NickName")),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`,`nick_name`,`balance`) VALUES (?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `nick_name`=VALUES(`nick_name`);",
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100),
				},
			},
		}, {
			name: "update column with assign field on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts().Update(ColWithUpdate("NickName", "Name")),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`,`nick_name`,`balance`) VALUES (?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `nick_name`=VALUES(`name`);",
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100),
				},
			},
		}, {
			name: "update column with value on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts().Update(Assign("NickName", "nick foo"), Assign("Balance", int64(10000))),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`,`nick_name`,`balance`) VALUES (?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `nick_name`=?,`balance`=?;",
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100), "nick foo", int64(10000),
				},
			},
		}, {
			name: "update column with value on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts().Update(invalidAssign{}),
			wantErr: errs.ErrInvalidAssignment,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stat, err := tc.inserter.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantStat, stat)
			}
		})
	}
}

func TestInserter_Build_StandardSQL(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, DBWithDialect(StandardSQL))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		inserter StatBuilder
		wantStat *Statement
		wantErr  error
	}{
		{
			name: "update column on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts("NickName").Update(Col("NickName")),
			wantStat: &Statement{
				SQL: `INSERT INTO "inserter_build_arg"("id","name","nick_name","balance") VALUES (?,?,?,?) ` +
					`ON CONFLICT ("nick_name") DO UPDATE SET "nick_name"=EXCLUDED."nick_name";`,
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100),
				},
			},
		}, {
			name: "update column with assign field on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts("NickName").Update(ColWithUpdate("NickName", "Name")),
			wantStat: &Statement{
				SQL: `INSERT INTO "inserter_build_arg"("id","name","nick_name","balance") VALUES (?,?,?,?) ` +
					`ON CONFLICT ("nick_name") DO UPDATE SET "nick_name"=EXCLUDED."name";`,
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100),
				},
			},
		}, {
			name: "update multi column with assign field on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts("Id").Update(Col("NickName"), Col("Balance")),
			wantStat: &Statement{
				SQL: `INSERT INTO "inserter_build_arg"("id","name","nick_name","balance") VALUES (?,?,?,?) ` +
					`ON CONFLICT ("id") DO UPDATE SET "nick_name"=EXCLUDED."nick_name","balance"=EXCLUDED."balance";`,
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100),
				},
			},
		}, {
			name: "update column with value on duplicate key",
			inserter: NewInserter[inserterBuildArg](db).Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
					NickName: &sql.NullString{
						Valid:  true,
						String: "foo bar",
					},
					Balance: int64(100),
				},
			).OnConflicts("Id", "Name").Update(Assign("NickName", "nick foo"), Assign("Balance", int64(10000))),
			wantStat: &Statement{
				SQL: `INSERT INTO "inserter_build_arg"("id","name","nick_name","balance") VALUES (?,?,?,?) ` +
					`ON CONFLICT ("id","name") DO UPDATE SET "nick_name"=?,"balance"=?;`,
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100), "nick foo", int64(10000),
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stat, err := tc.inserter.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantStat, stat)
			}
		})
	}
}

func TestInserter_Exec(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	tcs := []struct {
		name         string
		inserter     *Inserter[inserterBuildArg]
		wantErr      error
		rowsAffected int64
	}{
		{
			name: "invalid column",
			inserter: func() *Inserter[inserterBuildArg] {
				return NewInserter[inserterBuildArg](db).Row(&inserterBuildArg{}).ColFd("invalid")
			}(),
			wantErr: errs.ErrInvalidColumnFd("invalid"),
		}, {
			name: "db error",
			inserter: func() *Inserter[inserterBuildArg] {
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("mock db error"))

				return NewInserter[inserterBuildArg](db).Row(&inserterBuildArg{})
			}(),
			wantErr: errors.New("mock db error"),
		}, {
			name: "normal",
			inserter: func() *Inserter[inserterBuildArg] {
				mock.ExpectExec("INSERT INTO .*").
					WillReturnResult(driver.RowsAffected(1))

				return NewInserter[inserterBuildArg](db).Row(&inserterBuildArg{})
			}(),
			rowsAffected: int64(1),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.inserter.Exec(context.Background())
			assert.Equal(t, tc.wantErr, res.err)

			if res.err == nil {
				assert.Equal(t, tc.rowsAffected, res.RowsAffected())
			}
		})
	}
}
