package safedel

import (
	"context"
	"strings"

	easyorm "github.com/JrMarcco/easy-orm"
	"github.com/JrMarcco/easy-orm/internal/errs"
)

type MiddlewareBuilder struct{}

func (m *MiddlewareBuilder) Build() easyorm.Middleware {
	return func(next easyorm.HandleFunc) easyorm.HandleFunc {
		return func(ctx context.Context, ormCtx *easyorm.OrmContext) *easyorm.OrmResult {
			if ormCtx.Typ != easyorm.ScTypDELETE {
				return next(ctx, ormCtx)
			}

			statement, err := ormCtx.Builder.Build()
			if err != nil {
				return &easyorm.OrmResult{Err: err}
			}

			if strings.Contains(statement.SQL, "WHERE") {
				return &easyorm.OrmResult{
					Err: errs.ErrUnsafeDelete,
				}
			}

			return next(ctx, ormCtx)
		}
	}
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}
