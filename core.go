package easyorm

import (
	"context"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/internal/value"
	"github.com/JrMarcco/easy-orm/model"
)

type Core struct {
	dialect         Dialect
	registry        model.Registry
	resolverCreator value.ResolverCreator

	mws []Middleware
}

func findOneHF[T any](c context.Context, statementCtx *StatementContext, session session) *StatementResult {
	statement, err := statementCtx.Builder.Build()
	if err != nil {
		return &StatementResult{
			Err: err,
		}
	}

	rows, err := session.queryContext(c, statement.Sql, statement.Args...)
	if err != nil {
		return &StatementResult{
			Err: err,
		}
	}

	if !rows.Next() {
		return &StatementResult{
			Err: errs.ErrEligibleRow,
		}
	}

	res := new(T)

	m, err := session.getCore().registry.GetModel(res)
	if err != nil {
		return &StatementResult{
			Err: err,
		}
	}

	resolver := session.getCore().resolverCreator(m, res)
	if err = resolver.WriteColumns(rows); err != nil {
		return &StatementResult{
			Err: err,
		}
	}

	return &StatementResult{
		Res: res,
	}
}

func findOne[T any](c context.Context, sc *StatementContext, session session) (*T, error) {
	handleFunc := func(c context.Context, innerSc *StatementContext) *StatementResult {
		return findOneHF[T](c, innerSc, session)
	}

	core := session.getCore()
	for i := len(core.mws) - 1; i >= 0; i-- {
		handleFunc = core.mws[i](handleFunc)
	}

	sr := handleFunc(c, sc)
	if sr.Err != nil {
		return nil, sr.Err
	}
	return sr.Res.(*T), nil
}

func findMultiHF[T any](c context.Context, statementCtx *StatementContext, session session) *StatementResult {
	statement, err := statementCtx.Builder.Build()
	if err != nil {
		return &StatementResult{
			Err: err,
		}
	}

	rows, err := session.queryContext(c, statement.Sql, statement.Args...)
	if err != nil {
		return &StatementResult{
			Err: err,
		}
	}

	m, err := session.getCore().registry.GetModel(new(T))
	if err != nil {
		return &StatementResult{
			Err: err,
		}
	}

	res := make([]*T, 0, 16)
	for rows.Next() {
		v := new(T)

		resolver := session.getCore().resolverCreator(m, v)
		if err = resolver.WriteColumns(rows); err != nil {
			return &StatementResult{
				Err: err,
			}
		}

		res = append(res, v)
	}

	return &StatementResult{
		Res: res,
	}
}

func findMulti[T any](c context.Context, sc *StatementContext, session session) ([]*T, error) {
	handleFunc := func(c context.Context, innerSc *StatementContext) *StatementResult {
		return findMultiHF[T](c, innerSc, session)
	}

	core := session.getCore()
	for i := len(core.mws) - 1; i >= 0; i-- {
		handleFunc = core.mws[i](handleFunc)
	}

	sr := handleFunc(c, sc)
	if sr.Err != nil {
		return nil, sr.Err
	}
	return sr.Res.([]*T), nil
}
