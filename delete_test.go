package orm

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleter_Build(t *testing.T) {

	db, err := OpenDB(&sql.DB{}, DBWithDialect(MySqlDialect))
	require.NoError(t, err)

	tcs := []struct {
		name     string
		builder  StatBuilder
		wantStat *Statement
		wantErr  error
	}{
		{
			name:    "basic * select without from",
			builder: NewDeleter[deleterBuildArg](db),
			wantStat: &Statement{
				SQL: "DELETE FROM `deleter_build_arg`;",
			},
		},
		{
			name:    "basic * select with from",
			builder: NewDeleter[deleterBuildArg](db).From("test_model"),
			wantStat: &Statement{
				SQL: "DELETE FROM `test_model`;",
			},
		}, {
			name:    "basic * select with empty from",
			builder: NewDeleter[deleterBuildArg](db).From(""),
			wantStat: &Statement{
				SQL: "DELETE FROM `deleter_build_arg`;",
			},
		}, {
			name:    "basic * select with from db fdName",
			builder: NewDeleter[deleterBuildArg](db).From("test_db.test_model"),
			wantStat: &Statement{
				SQL: "DELETE FROM `test_db`.`test_model`;",
			},
		}, {
			name:    "empty where",
			builder: NewDeleter[deleterBuildArg](db).Where(),
			wantStat: &Statement{
				SQL: "DELETE FROM `deleter_build_arg`;",
			},
		}, {
			name:    "single predicate where",
			builder: NewDeleter[deleterBuildArg](db).Where(Col("Age").Eq(18)),
			wantStat: &Statement{
				SQL:  "DELETE FROM `deleter_build_arg` WHERE `age` = ?;",
				Args: []any{18},
			},
		}, {
			name:    "not predicate where",
			builder: NewDeleter[deleterBuildArg](db).Where(Not(Col("Age").Eq(18))),
			wantStat: &Statement{
				SQL:  "DELETE FROM `deleter_build_arg` WHERE NOT (`age` = ?);",
				Args: []any{18},
			},
		}, {
			name: "not & and predicate where",
			builder: NewDeleter[deleterBuildArg](db).Where(
				Not(
					Col("Age").Eq(18).And(Col("Id").Eq(1)),
				),
			),
			wantStat: &Statement{
				SQL:  "DELETE FROM `deleter_build_arg` WHERE NOT ((`age` = ?) AND (`id` = ?));",
				Args: []any{18, 1},
			},
		}, {
			name: "not & or predicate where",
			builder: NewDeleter[deleterBuildArg](db).Where(
				Not(
					Col("Id").Gt(100).Or(Col("Age").Lt(18)),
				),
			),
			wantStat: &Statement{
				SQL:  "DELETE FROM `deleter_build_arg` WHERE NOT ((`id` > ?) OR (`age` < ?));",
				Args: []any{100, 18},
			},
		}, {
			name:    "invalid type",
			builder: NewDeleter[deleterBuildArg](db).Where(Col("Invalid").Eq("test")),
			wantErr: errs.InvalidColumnFdErr("Invalid"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			stat, err := tc.builder.Build()
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantStat, stat)
			}
		})
	}

}

type deleterBuildArg struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  *sql.NullString
}
