package statlog

import (
	"context"
	orm "github.com/jrmarcco/easy-orm"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(stat string, args []any)
}

type MdlOpt func(builder *MiddlewareBuilder)

func BuilderWithLogFunc(logFunc func(stat string, args []any)) MdlOpt {
	return func(builder *MiddlewareBuilder) {
		builder.logFunc = logFunc
	}
}

func NewBuilder(opts ...MdlOpt) *MiddlewareBuilder {
	builder := &MiddlewareBuilder{
		logFunc: func(stat string, args []any) {
			log.Printf("statement: %s, args: %v \n", stat, args)
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

			sb, err := sc.Sb.Build()
			if err != nil {
				return &orm.StatResult{Err: err}
			}

			m.logFunc(sb.SQL, sb.Args)

			return next(ctx, sc)
		}
	}
}
