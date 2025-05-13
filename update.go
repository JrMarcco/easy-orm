package easyorm

import "context"

var _ Executor[any] = (*Updater[any])(nil)

type Updater[T any] struct {
	builder

	orm orm
}

func (u Updater[T]) Exec(ctx context.Context) Result {
	//TODO implement me
	panic("implement me")
}

func NewUpdater[T any](session orm) *Updater[T] {
	return &Updater[T]{
		builder: newBuilder(session),
		orm:     session,
	}
}
