package orm

import (
	"context"
	"database/sql"
)

// Querier 查询器
// select 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor 执行器
// insert / update / delete 语句
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

type Statement struct {
	SQL  string
	Args []any
}

type StatBuilder interface {
	Build() (*Statement, error)
}

type TbName interface {
	TbName() string
}
