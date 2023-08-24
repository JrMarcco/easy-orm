package otel

import (
	"context"
	orm "github.com/jrmarcco/easy-orm"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareBuilder struct {
	tracer trace.Tracer
}

func NewBuilder(opts ...MdlOpt) *MiddlewareBuilder {
	tracer := otel.GetTracerProvider().Tracer("github.com/jrmarcco/easy-orm/middleware/otel")

	builder := &MiddlewareBuilder{
		tracer: tracer,
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

type MdlOpt func(builder *MiddlewareBuilder)

func BuilderWithTracer(tracer trace.Tracer) MdlOpt {
	return func(builder *MiddlewareBuilder) {
		builder.tracer = tracer
	}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, sc *orm.StatContext) *orm.StatResult {

			tb := sc.Model.Tb
			extCtx, span := m.tracer.Start(ctx, sc.Typ+":"+tb, trace.WithAttributes())
			defer span.End()

			span.SetAttributes(attribute.String("component", "orm"))
			stat, err := sc.Builder.Build()
			if err != nil {
				span.RecordError(err)
			}
			span.SetAttributes(attribute.String("table", tb))

			if stat != nil {
				span.SetAttributes(attribute.String("sql", stat.SQL))
			}

			return next(extCtx, sc)
		}
	}
}
