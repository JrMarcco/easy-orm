package easyorm

import (
	"context"
	"database/sql"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/internal/value"
	"github.com/JrMarcco/easy-orm/model"
)

type Core struct {
	dialect         Dialect
	registry        model.Registry
	resolverCreator value.ResolverCreator

	middlewareChain MiddlewareChain
}

func findOneHF[T any](ctx context.Context, statementCtx *StatementContext, session session) *StatementResult {
	statement, err := statementCtx.Builder.Build()
	if err != nil {
		return &StatementResult{Err: err}
	}

	rows, err := session.queryContext(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return &StatementResult{Err: err}
	}

	if !rows.Next() {
		return &StatementResult{Err: errs.ErrEligibleRow}
	}

	res := new(T)
	m, err := session.getCore().registry.GetModel(res)
	if err != nil {
		return &StatementResult{Err: err}
	}

	resolver := session.getCore().resolverCreator(m, res)
	if err = resolver.WriteColumns(rows); err != nil {
		return &StatementResult{Err: err}
	}

	return &StatementResult{Res: res}
}

func findOne[T any](ctx context.Context, statementCtx *StatementContext, session session) (*T, error) {
	handleFunc := func(innerCtx context.Context, innerStatementCtx *StatementContext) *StatementResult {
		return findOneHF[T](innerCtx, innerStatementCtx, session)
	}

	core := session.getCore()
	for i := len(core.middlewareChain) - 1; i >= 0; i-- {
		handleFunc = core.middlewareChain[i](handleFunc)
	}

	sr := handleFunc(ctx, statementCtx)
	if sr.Err != nil {
		return nil, sr.Err
	}
	return sr.Res.(*T), nil
}

func findMultiHF[T any](ctx context.Context, statementCtx *StatementContext, session session) *StatementResult {
	statement, err := statementCtx.Builder.Build()
	if err != nil {
		return &StatementResult{Err: err}
	}

	rows, err := session.queryContext(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return &StatementResult{Err: err}
	}

	m, err := session.getCore().registry.GetModel(new(T))
	if err != nil {
		return &StatementResult{Err: err}
	}

	res := make([]*T, 0, 16)
	for rows.Next() {
		v := new(T)

		resolver := session.getCore().resolverCreator(m, v)
		if err = resolver.WriteColumns(rows); err != nil {
			return &StatementResult{Err: err}
		}
		res = append(res, v)
	}

	return &StatementResult{Res: res}
}

func findMulti[T any](ctx context.Context, statementCtx *StatementContext, session session) ([]*T, error) {
	handleFunc := func(innerCtx context.Context, innerStatementCtx *StatementContext) *StatementResult {
		return findMultiHF[T](innerCtx, innerStatementCtx, session)
	}

	core := session.getCore()
	for i := len(core.middlewareChain) - 1; i >= 0; i-- {
		handleFunc = core.middlewareChain[i](handleFunc)
	}

	sr := handleFunc(ctx, statementCtx)
	if sr.Err != nil {
		return nil, sr.Err
	}
	return sr.Res.([]*T), nil
}

func execHF(ctx context.Context, statementCtx *StatementContext, session session) *StatementResult {
	statement, err := statementCtx.Builder.Build()
	if err != nil {
		return &StatementResult{Err: err}
	}

	res, err := session.execContext(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return &StatementResult{Err: err}
	}

	return &StatementResult{Res: res}
}

func exec(ctx context.Context, statementCtx *StatementContext, session session) Result {
	handleFunc := func(innerCtx context.Context, innerStatementCtx *StatementContext) *StatementResult {
		return execHF(innerCtx, innerStatementCtx, session)
	}

	core := session.getCore()
	for i := len(core.middlewareChain) - 1; i >= 0; i-- {
		handleFunc = core.middlewareChain[i](handleFunc)
	}

	sr := handleFunc(ctx, statementCtx)
	if sr.Res == nil {
		return Result{err: sr.Err}
	}

	return Result{
		res: sr.Res.(sql.Result),
	}
}
