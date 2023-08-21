package orm

import (
	"context"
	"github.com/jrmarcco/easy-orm/model"
)

var (
	ScTypSelect = "SELECT"
	ScTypDelete = "DELETE"
	ScTypUpdate = "UPDATE"
	ScTypInsert = "INSERT"
)

type StatContext struct {
	Typ     string
	Builder StatBuilder
	Model   *model.Model
}

type StatResult struct {
	Res any
	Err error
}

type HandleFunc func(ctx context.Context, sc *StatContext) *StatResult

type Middleware func(next HandleFunc) HandleFunc
