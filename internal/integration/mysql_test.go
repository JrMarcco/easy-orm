//go:build integration

package integration

import (
	"context"
	orm "github.com/jrmarcco/easy-orm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlTestSuite struct {
	Suite
}

func (m *MysqlTestSuite) TearDownTest() {
	res := orm.NewRawStat[any](m.db, "truncate table `simple_struct`").Exec(context.Background())
	require.NoError(m.T(), res.Err())
}

func (m *MysqlTestSuite) TestInsert() {

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

	t := m.T()
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			res := orm.NewInserter[simpleStruct](m.db).Row(tc.rows...).Exec(context.Background())
			assert.Equal(t, tc.wantErr, res.Err())

			if res.Err() == nil {
				assert.Equal(t, tc.wantRes, res.RowsAffected())

				rows, err := orm.NewSelector[simpleStruct](m.db).GetMulti(context.Background())
				require.NoError(t, err)
				assert.Equal(t, tc.rows, rows)

				delRes := orm.NewDeleter[simpleStruct](m.db).Where(orm.Col("Id").Gt(0)).Exec(context.Background())
				require.NoError(t, delRes.Err())
				assert.Equal(t, tc.wantRes, delRes.RowsAffected())
			}
		})
	}
}

func TestMysql(t *testing.T) {
	suite.Run(t, &MysqlTestSuite{
		Suite: Suite{
			driver:  "mysql",
			dsn:     "jrmarcco:passwd@tcp(172.28.32.1:3306)/integration_test",
			dialect: orm.MySqlDialect,
		},
	})
}
