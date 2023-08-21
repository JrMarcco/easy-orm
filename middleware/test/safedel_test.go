package test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	orm "github.com/jrmarcco/easy-orm"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/middleware/safedel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleter_Exec(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	defer func(mockDB *sql.DB) {
		_ = mockDB.Close()
	}(mockDB)

	db, err := orm.OpenDB(mockDB, orm.DBWithDialect(orm.MySqlDialect), orm.DBWithMdls(safedel.NewBuilder().Build()))
	require.NoError(t, err)

	tcs := []struct {
		name    string
		deleter *orm.Deleter[deleterBuildArg]
		wantRes int64
		wantErr error
	}{
		{
			name:    "unsafe delete",
			deleter: orm.NewDeleter[deleterBuildArg](db),
			wantErr: errs.UnsafeDeleteErr,
		}, {
			name: "safe delete",
			deleter: func() *orm.Deleter[deleterBuildArg] {
				mock.ExpectExec("DELETE FROM `deleter_build_arg` WHERE `age` = ?").
					WillReturnResult(driver.RowsAffected(1))
				return orm.NewDeleter[deleterBuildArg](db).Where(orm.Col("Age").Eq(18))
			}(),
			wantRes: int64(1),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.deleter.Exec(context.Background())

			require.NotNil(t, res)

			assert.Equal(t, tc.wantErr, res.Err())
			if res.Err() == nil {
				assert.Equal(t, tc.wantRes, res.RowsAffected())
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
