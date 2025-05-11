package easyorm

import "context"

var _ Executor[any] = (*Inserter[any])(nil)

type Inserter[T any] struct {
	builder

	orm orm
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	//TODO implement me
	panic("implement me")
}

func (i *Inserter[T]) Build() (*Statement, error) {
	return &Statement{
		SQL:  i.sqlBuffer.String(),
		Args: i.args,
	}, nil
}

func NewInserter[T any](orm orm) *Inserter[T] {
	return &Inserter[T]{
		builder: newBuilder(orm),
		orm:     orm,
	}
}
