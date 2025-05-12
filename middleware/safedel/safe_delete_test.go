package safedel

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	easyorm "github.com/JrMarcco/easy-orm"
	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type safeDelTestModel struct {
	Id   uint64
	Name string
}

func TestMiddlewareBuilder_Build(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	opt := easyorm.DBWithMiddlewareChain(easyorm.MiddlewareChain{
		NewMiddlewareBuilder().Build(),
	})

	db, err := easyorm.OpenDB(mockDB, easyorm.MySQLDialect, opt)
	require.NoError(t, err)

	tcs := []struct {
		name    string
		deleter *easyorm.Deleter[safeDelTestModel]
		wantRes int64
		wantErr error
	}{
		{
			name: "basic",
			deleter: func() *easyorm.Deleter[safeDelTestModel] {
				mock.ExpectExec("DELETE FROM `safe_del_test_model`.*").
					WillReturnResult(sqlmock.NewResult(1, 1))

				return easyorm.NewDeleter[safeDelTestModel](db).Where(easyorm.Col("Id").Eq(1))
			}(),
			wantRes: 1,
		}, {
			name: "unsafe del",
			deleter: func() *easyorm.Deleter[safeDelTestModel] {
				mock.ExpectExec("DELETE FROM `safe_del_test_model`.*").
					WillReturnResult(sqlmock.NewResult(1, 1))

				return easyorm.NewDeleter[safeDelTestModel](db)
			}(),
			wantErr: errs.ErrUnsafeDelete,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.deleter.Exec(context.Background())
			assert.Equal(t, tc.wantErr, res.Err())

			if res.Err() == nil {
				assert.Equal(t, tc.wantRes, res.RowsAffected())
			}
		})
	}
}
