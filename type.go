package easyorm

import "context"

// Querier used to execute select statement.
type Querier[T any] interface {
	FindOne(ctx context.Context) (*T, error)
	FindMulti(ctx context.Context) ([]*T, error)
}

// Executor sql executor for insert, update, delete statement.
type Executor interface {
	Exec(ctx context.Context) (Result, error)
}

// Statement sql statement, include params.
type Statement struct {
	Sql  string
	Args []any
}

type StatementBuilder interface {
	Build() (*Statement, error)
}
