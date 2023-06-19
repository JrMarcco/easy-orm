package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"reflect"
	"regexp"
	"strings"
)

type registry struct {
	models map[reflect.Type]*model
}

func newRegistry() *registry {
	return &registry{
		models: make(map[reflect.Type]*model, 64),
	}
}

func (r *registry) getModel(entity any) (*model, error) {
	m, ok := r.models[reflect.TypeOf(entity)]
	if !ok {
		var err error
		m, err = r.parseModel(entity)
		if err != nil {
			return nil, err
		}

		r.models[reflect.TypeOf(entity)] = m
	}

	return m, nil
}

// parseModel 解析 model。
//
// entity 只能是结构体或指向结构体的一级指针。
func (r *registry) parseModel(entity any) (*model, error) {

	typ := reflect.TypeOf(entity)

	if typ.Kind() != reflect.Struct {

		if typ.Kind() != reflect.Pointer {
			return nil, errs.InvalidTypeErr
		}

		typ = typ.Elem()

		if typ.Kind() != reflect.Struct {
			return nil, errs.InvalidTypeErr
		}

	}

	numField := typ.NumField()
	fds := make(map[string]field, numField)

	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fds[fd.Name] = field{
			colName: camelToUnderline(fd.Name),
		}
	}

	m := &model{
		tbName: camelToUnderline(typ.Name()),
		fds:    fds,
	}

	r.models[typ] = m

	return m, nil
}

type model struct {
	tbName string
	fds    map[string]field
}

type field struct {
	colName string
}

var (
	matchNonAlphaNumeric = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	matchFirstCap        = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap          = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func camelToUnderline(camel string) string {

	camel = matchNonAlphaNumeric.ReplaceAllString(camel, "_")
	camel = matchFirstCap.ReplaceAllString(camel, "${1}_${2}")
	camel = matchAllCap.ReplaceAllString(camel, "${1}_${2}")

	return strings.ToLower(camel)

}
