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
		return func(ctx context.Context, statementCtx *easyorm.StatementContext) *easyorm.StatementResult {
			if statementCtx.Typ != easyorm.ScTypDELETE {
				return next(ctx, statementCtx)
			}

			statement, err := statementCtx.Builder.Build()
			if err != nil {
				return &easyorm.StatementResult{Err: err}
			}

			if strings.Contains(statement.SQL, "WHERE") {
				return &easyorm.StatementResult{
					Err: errs.ErrUnsafeDelete,
				}
			}

			return next(ctx, statementCtx)
		}
	}
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}
