package slowquery

import (
	"context"
	"fmt"
	orm "github.com/jrmarcco/easy-orm"
	"time"
)

type MiddlewareBuilder struct {
	// 阈值，参考值 100ms
	threshold time.Duration
	logFunc   func(stat *orm.Statement)
}

func NewBuilder(threshold time.Duration) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		threshold: threshold,
		logFunc: func(stat *orm.Statement) {
			fmt.Printf("slow query sql: %s\n", stat.SQL)
		},
	}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.HandleFunc) orm.HandleFunc {
		return func(ctx context.Context, sc *orm.StatContext) *orm.StatResult {

			start := time.Now()

			defer func() {
				duration := time.Since(start)

				if duration < m.threshold {
					return
				}

				// 慢查询
				if stat, err := sc.Builder.Build(); err == nil {
					m.logFunc(stat)
				}
			}()

			return next(ctx, sc)
		}
	}
}
