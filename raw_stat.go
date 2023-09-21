package orm

import (
	"context"
)

// RawStat 原生 sql 语句
type RawStat[T any] struct {
	builder

	sql  string
	args []any

	session Session
}

var _ Querier[any] = new(RawStat[any])

func NewRawStat[T any](session Session, sql string, args ...any) *RawStat[T] {
	return &RawStat[T]{
		builder: newBuilder(session),
		sql:     sql,
		args:    args,
		session: session,
	}
}

func (r *RawStat[T]) Build() (*Statement, error) {
	return &Statement{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

func (r *RawStat[T]) Get(ctx context.Context) (*T, error) {
	return get[T](ctx, r.session, &StatContext{
		Typ:     ScTypRaw,
		Builder: r,
	})
}

func (r *RawStat[T]) GetMulti(ctx context.Context) ([]*T, error) {
	return getMulti[T](ctx, r.session, &StatContext{
		Typ:     ScTypRaw,
		Builder: r,
	})
}

func (r *RawStat[T]) Exec(ctx context.Context) Result {
	return exec(ctx, r.session, &StatContext{
		Typ:     ScTypRaw,
		Builder: r,
	})
}
