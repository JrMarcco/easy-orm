package easyorm

import "context"

var _ Querier[any] = (*Raw[any])(nil)
var _ Executor[any] = (*Raw[any])(nil)

type Raw[T any] struct {
	builder

	session orm
	sql     string
	args    []any
}

func (r *Raw[T]) FindOne(ctx context.Context) (*T, error) {
	return findOne[T](ctx, &StatementContext{
		Typ:     ScTypRaw,
		Builder: r,
	}, r.session)
}

func (r *Raw[T]) FindMulti(ctx context.Context) ([]*T, error) {
	return findMulti[T](ctx, &StatementContext{
		Typ:     ScTypRaw,
		Builder: r,
	}, r.session)
}

func (r *Raw[T]) Build() (*Statement, error) {
	return &Statement{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

func (r *Raw[T]) Exec(ctx context.Context) Result {
	//TODO implement me
	panic("implement me")
}

func NewRaw[T any](session orm, sql string, args ...any) *Raw[T] {
	return &Raw[T]{
		builder: newBuilder(session),
		session: session,
		sql:     sql,
		args:    args,
	}
}

var _ selectable = (*RawExpression)(nil)
var _ Expression = (*RawExpression)(nil)

type RawExpression struct {
	raw  string
	args []any
}

func (r RawExpression) selectable() {}
func (r RawExpression) expr()       {}

func RawAsPd(raw string, args ...any) Predicate {
	re := RawExpression{
		raw:  raw,
		args: args,
	}

	return Predicate{
		left: re,
	}
}
