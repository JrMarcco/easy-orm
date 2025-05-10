package easyorm

import (
	"context"
)

const (
	ScTypRaw    = "RAW"
	ScTypSELECT = "SELECT"
	ScTypINSERT = "INSERT"
	ScTypUPDATE = "UPDATE"
	ScTypDELETE = "DELETE"
)

type StatementContext struct {
	Typ     string
	Builder StatementBuilder
}

type StatementResult struct {
	Res any
	Err error
}

type HandleFunc func(ctx context.Context, statementCtx *StatementContext) *StatementResult

type Middleware func(next HandleFunc) HandleFunc
