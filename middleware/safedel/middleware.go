package safedel

import (
	"context"
	"errors"
	orm "github.com/jrmarcco/easy-orm"
	"strings"
)

type MiddlewareBuilder struct {
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, sc *orm.StatContext) *orm.StatResult {

			if sc.Typ != orm.ScTypDelete {
				return next(ctx, sc)
			}

			stat, err := sc.Builder.Build()
			if err != nil {
				return &orm.StatResult{Err: err}
			}

			if !strings.Contains(stat.SQL, "WHERE") {
				return &orm.StatResult{
					Err: errors.New("unsafe operation, delete without where"),
				}
			}

			return next(ctx, sc)
		}
	}
}
