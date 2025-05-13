package easyorm

import (
	"context"
	"database/sql"

	"github.com/JrMarcco/easy-orm/internal/errs"
	"github.com/JrMarcco/easy-orm/internal/value"
	"github.com/JrMarcco/easy-orm/model"
)

type core struct {
	dialect         Dialect
	registry        model.Registry
	resolverCreator value.ResolverCreator

	middlewareChain MiddlewareChain
}

func findOneHF[T any](ctx context.Context, ormCtx *OrmContext, orm orm) *OrmResult {
	statement, err := ormCtx.Builder.Build()
	if err != nil {
		return &OrmResult{Err: err}
	}

	rows, err := orm.queryContext(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return &OrmResult{Err: err}
	}

	if !rows.Next() {
		return &OrmResult{Err: errs.ErrEligibleRow}
	}

	res := new(T)
	m, err := orm.getCore().registry.GetModel(res)
	if err != nil {
		return &OrmResult{Err: err}
	}

	resolver := orm.getCore().resolverCreator(m, res)
	if err = resolver.WriteColumns(rows); err != nil {
		return &OrmResult{Err: err}
	}

	return &OrmResult{Res: res}
}

func findOne[T any](ctx context.Context, ormCtx *OrmContext, orm orm) (*T, error) {
	handleFunc := func(innerCtx context.Context, innerOrmCtx *OrmContext) *OrmResult {
		return findOneHF[T](innerCtx, innerOrmCtx, orm)
	}

	c := orm.getCore()
	for i := len(c.middlewareChain) - 1; i >= 0; i-- {
		handleFunc = c.middlewareChain[i](handleFunc)
	}

	sr := handleFunc(ctx, ormCtx)
	if sr.Err != nil {
		return nil, sr.Err
	}
	return sr.Res.(*T), nil
}

func findMultiHF[T any](ctx context.Context, ormCtx *OrmContext, orm orm) *OrmResult {
	statement, err := ormCtx.Builder.Build()
	if err != nil {
		return &OrmResult{Err: err}
	}

	rows, err := orm.queryContext(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return &OrmResult{Err: err}
	}

	m, err := orm.getCore().registry.GetModel(new(T))
	if err != nil {
		return &OrmResult{Err: err}
	}

	res := make([]*T, 0, 16)
	for rows.Next() {
		v := new(T)

		resolver := orm.getCore().resolverCreator(m, v)
		if err = resolver.WriteColumns(rows); err != nil {
			return &OrmResult{Err: err}
		}
		res = append(res, v)
	}

	return &OrmResult{Res: res}
}

func findMulti[T any](ctx context.Context, ormCtx *OrmContext, orm orm) ([]*T, error) {
	handleFunc := func(innerCtx context.Context, innerOrmCtx *OrmContext) *OrmResult {
		return findMultiHF[T](innerCtx, innerOrmCtx, orm)
	}

	c := orm.getCore()
	for i := len(c.middlewareChain) - 1; i >= 0; i-- {
		handleFunc = c.middlewareChain[i](handleFunc)
	}

	sr := handleFunc(ctx, ormCtx)
	if sr.Err != nil {
		return nil, sr.Err
	}
	return sr.Res.([]*T), nil
}

func execHF(ctx context.Context, statementCtx *OrmContext, orm orm) *OrmResult {
	statement, err := statementCtx.Builder.Build()
	if err != nil {
		return &OrmResult{Err: err}
	}

	res, err := orm.execContext(ctx, statement.SQL, statement.Args...)
	if err != nil {
		return &OrmResult{Err: err}
	}

	return &OrmResult{Res: res}
}

func exec(ctx context.Context, ormCtx *OrmContext, orm orm) Result {
	handleFunc := func(innerCtx context.Context, innerOrmCtx *OrmContext) *OrmResult {
		return execHF(innerCtx, innerOrmCtx, orm)
	}

	c := orm.getCore()
	for i := len(c.middlewareChain) - 1; i >= 0; i-- {
		handleFunc = c.middlewareChain[i](handleFunc)
	}

	sr := handleFunc(ctx, ormCtx)
	if sr.Res == nil {
		return Result{err: sr.Err}
	}

	return Result{
		res: sr.Res.(sql.Result),
	}
}

type OrmContext struct {
	Typ     string
	Model   *model.Model
	Builder StatementBuilder
}

type OrmResult struct {
	Res any
	Err error
}
