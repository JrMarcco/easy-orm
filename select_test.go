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

type selectTestModel struct {
	Id       uint64
	Name     string
	Age      int8
	NickName *sql.NullString
}

type fromModel struct{}

func TestSelector_Build(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, MySQLDialect)
	require.NoError(t, err)

	tcs := []struct {
		name          string
		selector      StatementBuilder
		wantStatement *Statement
		wantErr       error
	}{
		{
			name:     "basic select",
			selector: NewSelector[selectTestModel](db),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model`;",
			},
			wantErr: nil,
		}, {
			name:     "with where",
			selector: NewSelector[selectTestModel](db).Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "with where not",
			selector: NewSelector[selectTestModel](db).Where(Col("Id").Eq(1).Not()),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE NOT (`id` = ?);",
				Args: []any{
					1,
				},
			},
		}, {
			name: "with where and",
			selector: NewSelector[selectTestModel](db).Where(
				Col("Id").Eq(1),
				Col("Name").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE (`id` = ?) AND (`name` = ?);",
				Args: []any{
					1, "jrmarcco",
				},
			},
		}, {
			name: "with where or",
			selector: NewSelector[selectTestModel](db).Where(
				Col("Age").Ge(18).
					Or(
						Col("Age").Lt(12),
					),
				Col("NickName").Eq("jrmarcco"),
			),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE ((`age` >= ?) OR (`age` < ?)) AND (`nick_name` = ?);",
				Args: []any{
					18, 12, "jrmarcco",
				},
			},
		}, {
			name: "with invalid column",
			selector: NewSelector[selectTestModel](db).Where(
				Col("InvalidColumn").Eq(1),
			),
			wantErr: errs.ErrInvalidField("InvalidColumn"),
		}, {
			name: "with basic selectable",
			selector: NewSelector[selectTestModel](db).
				Select(
					Col("Id"),
					Col("Name"),
					Col("Age"),
				).
				Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				SQL: "SELECT `id`, `name`, `age` FROM `select_test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name: "with tableAlias selectable",
			selector: NewSelector[selectTestModel](db).
				Select(
					Col("Id"),
					Col("Name").As("user_name"),
					Col("Age"),
				).
				Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				SQL: "SELECT `id`, `name` AS `user_name`, `age` FROM `select_test_model` WHERE `id` = ?;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "with count aggregate",
			selector: NewSelector[selectTestModel](db).Select(Count("Id")),
			wantStatement: &Statement{
				SQL: "SELECT COUNT(`id`) FROM `select_test_model`;",
			},
		}, {
			name:     "with sum aggregate",
			selector: NewSelector[selectTestModel](db).Select(Sum("Age")),
			wantStatement: &Statement{
				SQL: "SELECT SUM(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with max aggregate",
			selector: NewSelector[selectTestModel](db).Select(Max("Age")),
			wantStatement: &Statement{
				SQL: "SELECT MAX(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with min aggregate",
			selector: NewSelector[selectTestModel](db).Select(Min("Age")),
			wantStatement: &Statement{
				SQL: "SELECT MIN(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with avg aggregate",
			selector: NewSelector[selectTestModel](db).Select(Avg("Age")),
			wantStatement: &Statement{
				SQL: "SELECT AVG(`age`) FROM `select_test_model`;",
			},
		}, {
			name:     "with tableAlias aggregate",
			selector: NewSelector[selectTestModel](db).Select(Max("Age").As("max_age")),
			wantStatement: &Statement{
				SQL: "SELECT MAX(`age`) AS `max_age` FROM `select_test_model`;",
			},
		}, {
			name: "with invalid aggregate",
			selector: NewSelector[selectTestModel](db).Select(
				Max("InvalidColumn"),
			),
			wantErr: errs.ErrInvalidField("InvalidColumn"),
		}, {
			name: "with multiple aggregates",
			selector: NewSelector[selectTestModel](db).Select(
				Count("Id"),
				Sum("Age"),
				Max("Age"),
				Min("Age"),
				Avg("Age"),
			),
			wantStatement: &Statement{
				SQL: "SELECT COUNT(`id`), SUM(`age`), MAX(`age`), MIN(`age`), AVG(`age`) FROM `select_test_model`;",
			},
		}, {
			name: "with raw expression in where",
			selector: NewSelector[selectTestModel](db).
				Where(
					RawAsPd("`id` = (`age` + ?)", 1),
				),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` = (`age` + ?);",
				Args: []any{
					1,
				},
			},
			wantErr: nil,
		}, {
			name: "with raw expression in predicate",
			selector: NewSelector[selectTestModel](db).
				Where(
					Col("Id").Eq(RawAsPd("`id` + ?", 1)),
				),
			wantStatement: &Statement{
				SQL:  "SELECT * FROM `select_test_model` WHERE `id` = (`id` + ?);",
				Args: []any{1},
			},
		}, {
			name:     "with limit",
			selector: NewSelector[selectTestModel](db).Limit(1),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` LIMIT 1;",
			},
		}, {
			name:     "with offset",
			selector: NewSelector[selectTestModel](db).Offset(10),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` OFFSET 10;",
			},
		}, {
			name:     "with limit and offset",
			selector: NewSelector[selectTestModel](db).Limit(1).Offset(10),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` LIMIT 1 OFFSET 10;",
			},
		}, {
			name: "with group by",
			selector: NewSelector[selectTestModel](db).
				GroupBy(Col("Id"), Col("Age")).
				Where(Col("Id").Eq(1)),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` = ? GROUP BY `id`, `age`;",
				Args: []any{
					1,
				},
			},
		}, {
			name:     "with having without group by",
			selector: NewSelector[selectTestModel](db).Having(Col("Id").Eq(1)),
			wantErr:  errs.ErrHavingWithoutGroupBy,
		}, {
			name: "with having",
			selector: NewSelector[selectTestModel](db).
				Where(Col("Id").Eq(1)).
				GroupBy(Col("Age")).
				Having(Col("Age").Eq(18)),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` = ? GROUP BY `age` HAVING `age` = ?;",
				Args: []any{
					1, 18,
				},
			},
		}, {
			name:     "with single asc",
			selector: NewSelector[selectTestModel](db).OrderBy(Asc("Id")),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` ORDER BY `id` ASC;",
			},
		}, {
			name:     "with single desc",
			selector: NewSelector[selectTestModel](db).OrderBy(Desc("Id")),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` ORDER BY `id` DESC;",
			},
		}, {
			name:     "with multiple asc",
			selector: NewSelector[selectTestModel](db).OrderBy(Asc("Id"), Asc("Age")),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` ORDER BY `id` ASC, `age` ASC;",
			},
		}, {
			name:     "with multiple desc",
			selector: NewSelector[selectTestModel](db).OrderBy(Desc("Id"), Desc("Age")),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` ORDER BY `id` DESC, `age` DESC;",
			},
		}, {
			name: "with single asc and desc",
			selector: NewSelector[selectTestModel](db).OrderBy(
				Asc("Id"),
				Desc("Age"),
			),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` ORDER BY `id` ASC, `age` DESC;",
			},
		}, {
			name: "with multiple asc and desc",
			selector: NewSelector[selectTestModel](db).OrderBy(
				Asc("Id"),
				Desc("Age"),
				Asc("Name"),
			),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` ORDER BY `id` ASC, `age` DESC, `name` ASC;",
			},
		}, {
			name:     "with invalid order",
			selector: NewSelector[insertTestModel](db).OrderBy(Asc("InvalidColumn")),
			wantErr:  errs.ErrInvalidField("InvalidColumn"),
		}, {
			name:     "with from",
			selector: NewSelector[insertTestModel](db).From(TableOf(fromModel{})),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `from_model`;",
			},
		}, {
			name:     "with where in",
			selector: NewSelector[selectTestModel](db).Where(Col("Id").In(1, 2, 3)),
			wantStatement: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` IN (?,?,?);",
				Args: []any{
					1, 2, 3,
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			statement, err := tc.selector.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantStatement, statement)
			}
		})
	}
}

type firstModel struct {
	Id   uint64
	Name string

	UsingColFirst  int64
	UsingColSecond string
}

type secondModel struct {
	Id      uint64
	FirstId uint64
	ThirdId uint64

	UsingColFirst  int64
	UsingColSecond string
}

type thirdModel struct {
	Id uint64
}

func TestSelector_Join(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, PostgresDialect)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		selector *Selector[firstModel]
		wantRes  *Statement
		wantErr  error
	}{
		{
			name: "inner join with on",
			selector: func() *Selector[firstModel] {
				firstTable := TableOf(firstModel{})
				secondTable := TableOf(secondModel{})

				tableRef := firstTable.InnerJoin(secondTable).On(
					firstTable.Col("Id").Eq(secondTable.Col("FirstId")),
				)

				return NewSelector[firstModel](db).From(tableRef).Where(firstTable.Col("Id").Eq(1))
			}(),
			wantRes: &Statement{
				SQL:  `SELECT * FROM "first_model" INNER JOIN "second_model" ON "id" = "first_id" WHERE "id" = $1;`,
				Args: []any{1},
			},
		}, {
			name: "inner join with tableAlias and on",
			selector: func() *Selector[firstModel] {
				firstTable := TableAs(firstModel{}, "f")
				secondTable := TableAs(secondModel{}, "s")

				tableRef := firstTable.InnerJoin(secondTable).On(
					firstTable.Col("Id").Eq(secondTable.Col("FirstId")),
				)

				return NewSelector[firstModel](db).From(tableRef).Where(firstTable.Col("Id").Eq(1))
			}(),
			wantRes: &Statement{
				SQL:  `SELECT * FROM "first_model" AS "f" INNER JOIN "second_model" AS "s" ON "f"."id" = "s"."first_id" WHERE "f"."id" = $1;`,
				Args: []any{1},
			},
		}, {
			name: "left join with using single column",
			selector: func() *Selector[firstModel] {
				firstTable := TableOf(firstModel{})
				secondTable := TableOf(secondModel{})

				tableRef := firstTable.LeftJoin(secondTable).Using(Col("UsingColFirst"))

				return NewSelector[firstModel](db).From(tableRef).Where(firstTable.Col("Id").Eq(1))
			}(),
			wantRes: &Statement{
				SQL:  `SELECT * FROM "first_model" LEFT JOIN "second_model" USING ("using_col_first") WHERE "id" = $1;`,
				Args: []any{1},
			},
		}, {
			name: "left join with using multiple columns",
			selector: func() *Selector[firstModel] {
				firstTable := TableOf(firstModel{})
				secondTable := TableOf(secondModel{})

				tableRef := firstTable.LeftJoin(secondTable).Using(
					Col("UsingColFirst"),
					Col("UsingColSecond"),
				)

				return NewSelector[firstModel](db).From(tableRef).Where(firstTable.Col("Id").Eq(1))
			}(),
			wantRes: &Statement{
				SQL:  `SELECT * FROM "first_model" LEFT JOIN "second_model" USING ("using_col_first", "using_col_second") WHERE "id" = $1;`,
				Args: []any{1},
			},
		}, {
			name: "right join with using invalid column",
			selector: func() *Selector[firstModel] {
				firstTable := TableOf(firstModel{})
				secondTable := TableOf(secondModel{})

				tableRef := firstTable.LeftJoin(secondTable).Using(Col("Invalid"))
				return NewSelector[firstModel](db).From(tableRef).Where(firstTable.Col("Id").Eq(1))
			}(),
			wantErr: errs.ErrInvalidField("Invalid"),
		}, {
			name: "right join after inner join",
			selector: func() *Selector[firstModel] {
				firstTable := TableOf(firstModel{}).As("f")
				secondTable := TableOf(secondModel{}).As("s")
				thirdTable := TableOf(thirdModel{}).As("t")

				innerJoin := firstTable.InnerJoin(secondTable).On(firstTable.Col("Id").Eq(secondTable.Col("FirstId")))

				rightJoin := innerJoin.RightJoin(thirdTable).On(
					secondTable.Col("ThirdId").Eq(thirdTable.Col("Id")),
				)

				return NewSelector[firstModel](db).From(rightJoin).Where(firstTable.Col("Id").Gt(100))
			}(),
			wantRes: &Statement{
				SQL:  `SELECT * FROM "first_model" AS "f" INNER JOIN "second_model" AS "s" ON "f"."id" = "s"."first_id" RIGHT JOIN "third_model" AS "t" ON "s"."third_id" = "t"."id" WHERE "f"."id" > $1;`,
				Args: []any{100},
			},
		}, {
			name: "inner join after left join",
			selector: func() *Selector[firstModel] {
				firstTable := TableOf(firstModel{}).As("f")
				secondTable := TableOf(secondModel{}).As("s")
				thirdTable := TableOf(thirdModel{}).As("t")

				leftJoin := firstTable.LeftJoin(secondTable).Using(
					firstTable.Col("UsingColFirst"),
					firstTable.Col("UsingColSecond"),
				)
				innerJoin := leftJoin.InnerJoin(thirdTable).On(
					secondTable.Col("ThirdId").Eq(thirdTable.Col("Id")),
				)

				return NewSelector[firstModel](db).From(innerJoin).Where(firstTable.Col("Id").Gt(100))
			}(),
			wantRes: &Statement{
				SQL:  `SELECT * FROM "first_model" AS "f" LEFT JOIN "second_model" AS "s" USING ("using_col_first", "using_col_second") INNER JOIN "third_model" AS "t" ON "s"."third_id" = "t"."id" WHERE "f"."id" > $1;`,
				Args: []any{100},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			statement, err := tc.selector.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, statement)
			}
		})
	}
}

func TestSelector_SubQuery(t *testing.T) {
	db, err := OpenDB(&sql.DB{}, MySQLDialect)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		selector *Selector[selectTestModel]
		wantRes  *Statement
		wantErr  error
	}{
		{
			name: "select from sub query",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[selectTestModel](db).AsSubQuery("s")
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).From(subQuery)
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM (SELECT * FROM `select_test_model`) AS `s`;",
			},
		}, {
			name: "select from sub query with where",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[selectTestModel](db).AsSubQuery("s")
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).From(subQuery).Where(subQuery.Col("Id").Eq(1))
			}(),
			wantRes: &Statement{
				SQL:  "SELECT * FROM (SELECT * FROM `select_test_model`) AS `s` WHERE `s`.`id` = ?;",
				Args: []any{1},
			},
		}, {
			name: "select from sub query with in",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[selectTestModel](db).Select(Col("Id")).ToSubQuery()
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).Where(Col("Id").InSubQuery(subQuery))
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` IN (SELECT `id` FROM `select_test_model`);",
			},
		}, {
			name: "where exists",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[secondModel](db).Select(Col("FirstId")).ToSubQuery()
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).Where(subQuery.Exists())
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE EXISTS (SELECT `first_id` FROM `second_model`);",
			},
		}, {
			name: "where not exists",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[secondModel](db).Select(Col("FirstId")).ToSubQuery()
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).Where(subQuery.NotExists())
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE NOT EXISTS (SELECT `first_id` FROM `second_model`);",
			},
		}, {
			name: "where all",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[secondModel](db).Select(Col("FirstId")).ToSubQuery()
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).Where(Col("Id").Gt(subQuery.All()))
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` > (ALL (SELECT `first_id` FROM `second_model`));",
			},
		}, {
			name: "where any",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[secondModel](db).Select(Col("FirstId")).ToSubQuery()
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).Where(Col("Id").Gt(subQuery.Any()))
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` > (ANY (SELECT `first_id` FROM `second_model`));",
			},
		}, {
			name: "where some",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[secondModel](db).Select(Col("FirstId")).ToSubQuery()
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).Where(Col("Id").Gt(subQuery.Some()))
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` > (SOME (SELECT `first_id` FROM `second_model`));",
			},
		}, {
			name: "where some and any",
			selector: func() *Selector[selectTestModel] {
				subQuery, err := NewSelector[secondModel](db).Select(Col("FirstId")).ToSubQuery()
				require.NoError(t, err)
				return NewSelector[selectTestModel](db).Where(Col("Id").Gt(subQuery.Some().And(subQuery.Any())))
			}(),
			wantRes: &Statement{
				SQL: "SELECT * FROM `select_test_model` WHERE `id` > ((SOME (SELECT `first_id` FROM `second_model`)) AND (ANY (SELECT `first_id` FROM `second_model`)));",
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			statement, err := tc.selector.Build()
			assert.Equal(t, tc.wantErr, err)

			if err == nil {
				assert.Equal(t, tc.wantRes, statement)
			}
		})
	}
}

func TestSelector_FindOne(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, MySQLDialect)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *Selector[selectTestModel]
		wantRes  *selectTestModel
		wantErr  error
	}{
		{
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
				rows.AddRow(1, "foo", 18, "bar")

				mock.ExpectQuery("SELECT *.").WillReturnRows(rows)
			},
			selector: NewSelector[selectTestModel](db),
			wantRes: &selectTestModel{
				Id:       1,
				Name:     "foo",
				Age:      18,
				NickName: &sql.NullString{String: "bar", Valid: true},
			},
		}, {
			name: "with args",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
				rows.AddRow(1, "foo", 18, "bar")

				mock.ExpectQuery("SELECT *.").WithArgs(1, "foo").WillReturnRows(rows)
			},
			selector: NewSelector[selectTestModel](db).
				Where(
					Col("Id").Eq(1),
					Col("Name").Eq("foo"),
				),
			wantRes: &selectTestModel{
				Id:       1,
				Name:     "foo",
				Age:      18,
				NickName: &sql.NullString{String: "bar", Valid: true},
			},
		}, {
			name: "returns error",
			mockFunc: func() {
				mock.ExpectQuery("SELECT *.").WillReturnError(errors.New("mock error"))
			},
			selector: NewSelector[selectTestModel](db),
			wantErr:  errors.New("mock error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			var res *selectTestModel
			res, err = tc.selector.FindOne(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}

func TestSelector_FindMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() {
		_ = mockDB.Close()
	}()

	db, err := OpenDB(mockDB, MySQLDialect)
	require.NoError(t, err)

	tcs := []struct {
		name     string
		mockFunc func()
		selector *Selector[selectTestModel]
		wantRes  []*selectTestModel
		wantErr  error
	}{
		{
			name: "basic",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "age", "nick_name"})
				rows.AddRow(1, "foo", 18, "bar")
				rows.AddRow(2, "bar", 12, "foo")
				rows.AddRow(3, "baz", 10, "baz")

				mock.ExpectQuery("SELECT *.").WillReturnRows(rows)
			},
			selector: NewSelector[selectTestModel](db),
			wantRes: []*selectTestModel{
				{
					Id:       1,
					Name:     "foo",
					Age:      18,
					NickName: &sql.NullString{String: "bar", Valid: true},
				}, {
					Id:       2,
					Name:     "bar",
					Age:      12,
					NickName: &sql.NullString{String: "foo", Valid: true},
				}, {
					Id:       3,
					Name:     "baz",
					Age:      10,
					NickName: &sql.NullString{String: "baz", Valid: true},
				},
			},
		}, {
			name: "returns error",
			mockFunc: func() {
				mock.ExpectQuery("SELECT *.").WillReturnError(errors.New("mock error"))
			},
			selector: NewSelector[selectTestModel](db),
			wantErr:  errors.New("mock error"),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			var res []*selectTestModel
			res, err = tc.selector.FindMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err == nil {
				assert.Equal(t, tc.wantRes, res)
			}
		})
	}
}
