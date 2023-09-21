package orm

import (
	"context"
	"github.com/jrmarcco/easy-orm/model"
)

const (
	ScTypRaw    = "RAW"
	ScTypSelect = "SELECT"
	ScTypDelete = "DELETE"
	ScTypUpdate = "UPDATE"
	ScTypInsert = "INSERT"
)

// StatContext sql statement context
// include:
//
//	statement type
//	statement builder
//	related model
type StatContext struct {
	Typ     string
	Builder StatBuilder
	Model   *model.Model
}

// StatResult sql statement exec result
type StatResult struct {
	Res any
	Err error
}

type HandleFunc func(ctx context.Context, sc *StatContext) *StatResult

type Middleware func(next HandleFunc) HandleFunc
