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
			wantErr:  errs.EmptyInserRowErr,
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
			name: "assign insert column field",
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
			name: "assign invalid insert column field",
			inserter: NewInserter[inserterBuildArg](db).ColFd("Id", "Invalid").Row(
				&inserterBuildArg{
					Id:   uint64(1),
					Name: "jrmarcco",
				},
			),
			wantErr: errs.InvalidColumnFdErr("Invalid"),
		}, {
			name: "assign insert column field with multi row",
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
