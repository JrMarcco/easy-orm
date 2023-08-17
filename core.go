package orm

import (
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
