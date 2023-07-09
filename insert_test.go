package orm

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db, err := OpenDB(&sql.DB{})
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
			wantErr:  errs.EmptyInsertRowErr,
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
			wantErr: errs.InvalidColumnFdErr("Invalid"),
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
			).OnDuplicateKey().Update(Col("NickName")),
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
			).OnDuplicateKey().Update(ColWithUpdate("NickName", "Name")),
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
			).OnDuplicateKey().Update(Assign("NickName", "jrmarcco"), Assign("Balance", int64(10000))),
			wantStat: &Statement{
				SQL: "INSERT INTO `inserter_build_arg`(`id`,`name`,`nick_name`,`balance`) VALUES (?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `nick_name`=?,`balance`=?;",
				Args: []any{
					uint64(1), "jrmarcco", &sql.NullString{Valid: true, String: "foo bar"}, int64(100), "jrmarcco", int64(10000),
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
			).OnDuplicateKey().Update(invalidAssign{}),
			wantErr: errs.InvalidAssignmentErr,
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

type inserterBuildArg struct {
	Id       uint64
	Name     string
	NickName *sql.NullString
	Balance  int64
}

type invalidAssign struct {
}

func (i invalidAssign) assign() {}
