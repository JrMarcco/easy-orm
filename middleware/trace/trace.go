package trace

import (
	"context"

	easyorm "github.com/JrMarcco/easy-orm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const defaultInstrumentationName = "github.com/JrMarcco/easy-orm/middleware/trace"

type MiddlewareBuilder struct {
	tracer trace.Tracer
}

func (m *MiddlewareBuilder) WithTracer(tracer trace.Tracer) *MiddlewareBuilder {
	m.tracer = tracer
	return m
}

func (m *MiddlewareBuilder) Build() easyorm.Middleware {
	if m.tracer == nil {
		m.tracer = otel.GetTracerProvider().Tracer(defaultInstrumentationName)
	}

	return func(next easyorm.HandleFunc) easyorm.HandleFunc {
		return func(ctx context.Context, ormCtx *easyorm.OrmContext) *easyorm.OrmResult {
			tableName := ormCtx.Model.TableName
			spanCtx, span := m.tracer.Start(ctx, tableName, trace.WithSpanKind(trace.SpanKindClient))
			defer span.End()

			span.SetAttributes(attribute.String("orm.type", ormCtx.Typ))
			span.SetAttributes(attribute.String("orm.table", tableName))

			statement, err := ormCtx.Builder.Build()
			if err != nil {
				span.RecordError(err)
				return &easyorm.OrmResult{Err: err}
			}

			if statement != nil {
				span.SetAttributes(attribute.String("orm.sql", statement.SQL))
			}

			res := next(spanCtx, ormCtx)
			if res.Err != nil {
				span.RecordError(res.Err)
			}
			return res
		}
	}
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}
