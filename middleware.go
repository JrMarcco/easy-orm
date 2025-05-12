package easyorm

import "context"

type HandleFunc func(ctx context.Context, ormCtx *OrmContext) *OrmResult

type Middleware func(next HandleFunc) HandleFunc

type MiddlewareChain []Middleware

const (
	ScTypRaw    = "RAW"
	ScTypSELECT = "SELECT"
	ScTypINSERT = "INSERT"
	ScTypUPDATE = "UPDATE"
	ScTypDELETE = "DELETE"
)
