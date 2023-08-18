package orm

import "context"

var (
	ScTypSelect = "SELECT"
	ScTypDelete = "DELETE"
	ScTypUpdate = "UPDATE"
	ScTypInsert = "INSERT"
)

type StatContext struct {
	Typ string
	Sb  StatBuilder
}

type StatResult struct {
	Res any
	Err error
}

type HandleFunc func(ctx context.Context, sc *StatContext) *StatResult

type Middleware func(next HandleFunc) HandleFunc