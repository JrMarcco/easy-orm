package orm

import (
	"context"
	"github.com/jrmarcco/easy-orm/internal/errs"
	"github.com/jrmarcco/easy-orm/internal/val"
	"github.com/jrmarcco/easy-orm/model"
)

type Core struct {
	registry model.Registry
	creator  val.Creator
	dialect  Dialect

	// AOP
	mdls []Middleware
}

func getSelectHF[T any](ctx context.Context, session Session, sc *StatContext) *StatResult {
	stat, err := sc.Builder.Build()
	if err != nil {
		return &StatResult{Err: err}
	}

	rows, err := session.queryContext(ctx, stat.SQL, stat.Args...)
	if err != nil {
		return &StatResult{Err: err}
	}

	if !rows.Next() {
		return &StatResult{Err: errs.NoEligibleRowsErr}
	}

	res := new(T)

	md, err := session.getCore().registry.Get(res)
	if err != nil {
		return &StatResult{Err: err}
	}

	writer := session.getCore().creator(md, res)
	if err = writer.WriteCols(rows); err != nil {
		return &StatResult{Err: err}
	}

	return &StatResult{Res: res}
}

func get[T any](ctx context.Context, session Session, sc *StatContext) (*T, error) {

	var hf = func(ctx context.Context, innerSc *StatContext) *StatResult {
		return getSelectHF[T](ctx, session, innerSc)
	}

	core := session.getCore()
	for i := len(core.mdls) - 1; i >= 0; i-- {
		hf = core.mdls[i](hf)
	}

	sr := hf(ctx, sc)

	if sr.Res != nil {
		return sr.Res.(*T), nil
	}
	return nil, sr.Err
}

func getSelectMultiHF[T any](ctx context.Context, session Session, sc *StatContext) *StatResult {
	stat, err := sc.Builder.Build()
	if err != nil {
		return &StatResult{Err: err}
	}

	rows, err := session.queryContext(ctx, stat.SQL, stat.Args...)
	if err != nil {
		return &StatResult{Err: err}
	}

	md, err := session.getCore().registry.Get(new(T))
	if err != nil {
		return &StatResult{Err: err}
	}

	res := make([]*T, 0, 8)
	for rows.Next() {
		value := new(T)

		writer := session.getCore().creator(md, value)
		if err := writer.WriteCols(rows); err != nil {
			return &StatResult{Err: err}
		}

		res = append(res, value)
	}
	return &StatResult{Res: res}
}

func getMulti[T any](ctx context.Context, session Session, sc *StatContext) ([]*T, error) {

	var hf = func(ctx context.Context, innerSc *StatContext) *StatResult {
		return getSelectMultiHF[T](ctx, session, innerSc)
	}

	core := session.getCore()
	for i := len(core.mdls) - 1; i >= 0; i-- {
		hf = core.mdls[i](hf)
	}

	sr := hf(ctx, sc)

	if sr.Res != nil {
		return sr.Res.([]*T), nil
	}

	return nil, sr.Err
}
