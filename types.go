package orm

import (
	"context"
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
	Exec(ctx context.Context) Result
}

// Statement sql语句，包含构造好的 sql 以及参数
type Statement struct {
	SQL  string
	Args []any
}

// StatBuilder sql语句构造器
type StatBuilder interface {
	Build() (*Statement, error)
}
