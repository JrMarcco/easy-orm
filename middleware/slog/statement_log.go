package slog

import (
	"context"
	"log"

	easyorm "github.com/JrMarcco/easy-orm"
)

type MiddlewareBuilder struct {
	logFunc func(s *easyorm.Statement)
}

type Opt func(*MiddlewareBuilder)

func WithLogFunc(logFunc func(s *easyorm.Statement)) Opt {
	return func(builder *MiddlewareBuilder) {
		builder.logFunc = logFunc
	}
}

func NewMiddlewareBuilder(opts ...Opt) *MiddlewareBuilder {
	builder := &MiddlewareBuilder{
		logFunc: func(s *easyorm.Statement) {
			log.Printf("[easy-orm] statement: [%s], args: [%v]\n", s.SQL, s.Args)
		},
	}
	for _, opt := range opts {
		opt(builder)
	}
	return builder
}

func (m *MiddlewareBuilder) Build() easyorm.Middleware {
	return func(next easyorm.HandleFunc) easyorm.HandleFunc {
		return func(ctx context.Context, ormCtx *easyorm.OrmContext) *easyorm.OrmResult {
			statement, err := ormCtx.Builder.Build()
			if err != nil {
				return &easyorm.OrmResult{Err: err}
			}

			m.logFunc(statement)

			return next(ctx, ormCtx)
		}
	}
}
