package integration

import (
	"context"
	orm "github.com/jrmarcco/easy-orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"

	_ "github.com/lib/pq"
)

type PostgresTestSuite struct {
	Suite
}

func (p *PostgresTestSuite) TestTearDown() {
	res := orm.NewRawStat[any](p.db, `-- truncate table "simple_struct"`).Exec(context.Background())
	require.NoError(p.T(), res.Err())
}

func (p *PostgresTestSuite) TestInsert() {

	boolVal := false
	int8Val := int8(8)
	int16Val := int16(16)

	tcs := []struct {
		name string

		rows    []*simpleStruct
		wantRes int64
		wantErr error
	}{
		{
			name: "single row",
			rows: []*simpleStruct{
				{
					Id: 1,
				},
			},
			wantRes: 1,
		}, {
			name: "single row with json",
			rows: []*simpleStruct{
				{
					Id: 2,
					JsonColumn: &simpleJson{
						Val:   simpleUser{Name: "jrmarcco"},
						Valid: true,
					},
				},
			},
			wantRes: 1,
		}, {
			name: "multi rows",
			rows: []*simpleStruct{
				{
					Id:      3,
					Bool:    true,
					BoolPtr: &boolVal,
				}, {
					Id:      4,
					Int8:    int8Val,
					Int8Ptr: &int8Val,
				}, {
					Id:       5,
					Int16:    int16Val,
					Int16Ptr: &int16Val,
				},
			},
			wantRes: 3,
		},
	}

	t := p.T()
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			res := orm.NewInserter[simpleStruct](p.db).Row(tc.rows...).Exec(context.Background())
			assert.Equal(t, tc.wantErr, res.Err())

			if res.Err() == nil {
				assert.Equal(t, tc.wantRes, res.RowsAffected())

				rows, err := orm.NewSelector[simpleStruct](p.db).GetMulti(context.Background())

				for _, row := range rows {
					// postgresql 在处理 ByteArray 字段时候
					// 会将 nil 作为空切片传入
					if len(row.ByteArray) == 0 {
						row.ByteArray = nil
					}
				}

				require.NoError(t, err)
				assert.Equal(t, tc.rows, rows)

				delRes := orm.NewDeleter[simpleStruct](p.db).Where(orm.Col("Id").Gt(0)).Exec(context.Background())
				require.NoError(t, delRes.Err())
				assert.Equal(t, tc.wantRes, delRes.RowsAffected())
			}
		})
	}
}

func TestPostgres(t *testing.T) {
	suite.Run(t, &PostgresTestSuite{
		Suite: Suite{
			driver:  "postgres",
			dsn:     "postgres://jrmarcco:example@172.28.32.1:5432/integration_test?sslmode=disable",
			dialect: orm.PostgresDialect,
		},
	})
}
