package prometheus

import (
	"context"
	"time"

	easyorm "github.com/JrMarcco/easy-orm"
	"github.com/prometheus/client_golang/prometheus"
)

type MiddlewareBuilder struct {
	vec *prometheus.SummaryVec
}

func (m *MiddlewareBuilder) WithSummaryVec(vec *prometheus.SummaryVec) *MiddlewareBuilder {
	m.vec = vec
	return m
}

func (m *MiddlewareBuilder) Build() easyorm.Middleware {
	if m.vec == nil {
		m.vec = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:        "easy_orm",
				Subsystem:   "sql_execute_duration",
				ConstLabels: map[string]string{},
				Objectives: map[float64]float64{
					0.5:   0.01,
					0.75:  0.01,
					0.90:  0.005,
					0.98:  0.002,
					0.99:  0.001,
					0.999: 0.0001,
				},
				Help: "Duration of SQL execute in microseconds.",
			},
			[]string{"type", "table"},
		)
	}

	return func(next easyorm.HandleFunc) easyorm.HandleFunc {
		return func(ctx context.Context, ormCtx *easyorm.OrmContext) *easyorm.OrmResult {
			start := time.Now()
			defer func() {
				m.vec.WithLabelValues(ormCtx.Typ, ormCtx.Model.TableName).Observe(float64(time.Since(start).Milliseconds()))
			}()
			return next(ctx, ormCtx)
		}
	}
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}
