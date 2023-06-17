package orm

import (
	"database/sql"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {

	tcs := []struct {
		name     string
		builder  StatBuilder
		wantStat *Statement
		wantErr  error
	}{
		{
			name:    "basic * select without from",
			builder: &Selector[selectorBuildArg]{},
			wantStat: &Statement{
				SQL: "SELECT * FROM `selector_build_arg`;",
			},
		}, {
			name:    "basic * select with from",
			builder: (&Selector[selectorBuildArg]{}).From("test_model"),
			wantStat: &Statement{
				SQL: "SELECT * FROM `test_model`;",
			},
		}, {
			name:    "basic * select with empty from",
			builder: (&Selector[selectorBuildArg]{}).From(""),
			wantStat: &Statement{
				SQL: "SELECT * FROM `selector_build_arg`;",
			},
		}, {
			name:    "basic * select with from db name",
			builder: (&Selector[selectorBuildArg]{}).From("test_db.test_model"),
			wantStat: &Statement{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		}, {
			name:    "empty where",
			builder: (&Selector[selectorBuildArg]{}).Where(),
			wantStat: &Statement{
				SQL: "SELECT * FROM `selector_build_arg`;",
			},
		}, {
			name:    "single predicate where",
			builder: (&Selector[selectorBuildArg]{}).Where(Col("Age").Eq(18)),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE `age` = ?;",
				Args: []any{18},
			},
		}, {
			name:    "not predicate where",
			builder: (&Selector[selectorBuildArg]{}).Where(Not(Col("Age").Eq(18))),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE NOT (`age` = ?);",
				Args: []any{18},
			},
		}, {
			name: "not & and predicate where",
			builder: (&Selector[selectorBuildArg]{}).Where(
				Not(
					Col("Age").Eq(18).And(Col("Id").Eq(1)),
				),
			),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE NOT ((`age` = ?) AND (`id` = ?));",
				Args: []any{18, 1},
			},
		}, {
			name: "not & or predicate where",
			builder: (&Selector[selectorBuildArg]{}).Where(
				Not(
					Col("Id").Gt(100).Or(Col("Age").Lt(18)),
				),
			),
			wantStat: &Statement{
				SQL:  "SELECT * FROM `selector_build_arg` WHERE NOT ((`id` > ?) OR (`age` < ?));",
				Args: []any{100, 18},
			},
		}, {
			name:    "invalid type",
			builder: (&Selector[selectorBuildArg]{}).Where(Col("Invalid").Eq("test")),
			wantErr: errs.InvalidColumnErr("Invalid"),
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

type selectorBuildArg struct {
	Id        int64
	Age       int8
	FirstName string
	LastName  *sql.NullString
}
