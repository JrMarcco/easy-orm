package statlog

import (
	"context"
	orm "github.com/jrmarcco/easy-orm"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(stat *orm.Statement)
}

type MdlOpt func(builder *MiddlewareBuilder)

func BuilderWithLogFunc(logFunc func(stat *orm.Statement)) MdlOpt {
	return func(builder *MiddlewareBuilder) {
		builder.logFunc = logFunc
	}
}

func NewBuilder(opts ...MdlOpt) *MiddlewareBuilder {
	builder := &MiddlewareBuilder{
		logFunc: func(stat *orm.Statement) {
			log.Printf("statement: %s, args: %v \n", stat.SQL, stat.Args)
		},
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, sc *orm.StatContext) *orm.StatResult {

			stat, err := sc.Sb.Build()
			if err != nil {
				return &orm.StatResult{Err: err}
			}

			m.logFunc(stat)

			return next(ctx, sc)
		}
	}
}
