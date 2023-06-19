package orm

import (
	"github.com/jrmarcco/easy-orm/internal/errs"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

const (
	tagKeyCol = "column"
)

type registry struct {
	sync.RWMutex
	models map[reflect.Type]*model
}

func newRegistry() *registry {
	return &registry{
		models: make(map[reflect.Type]*model, 64),
	}
}

func (r *registry) getModel(entity any) (*model, error) {

	typ := reflect.TypeOf(entity)

	r.RLock()
	m, ok := r.models[typ]
	r.RUnlock()

	if ok {
		return m, nil
	}

	r.Lock()
	defer r.Unlock()

	// double check
	if m, ok = r.models[typ]; ok {
		return m, nil
	}

	var err error
	m, err = r.parseModel(entity)
	if err != nil {
		return nil, err
	}

	r.models[typ] = m

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

		tags, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}

		colName, ok := tags[tagKeyCol]
		if !ok {
			colName = camelToUnderline(fd.Name)
		}

		fds[fd.Name] = field{
			colName: colName,
		}
	}

	var tbName string
	if ti, ok := entity.(TbName); ok {
		tbName = ti.TbName()
	}

	if tbName == "" {
		tbName = camelToUnderline(typ.Name())
	}

	m := &model{
		tbName: tbName,
		fds:    fds,
	}

	return m, nil
}

// parseTag 解析标签
func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("orm")
	if ormTag == "" {
		return map[string]string{}, nil
	}

	pairs := strings.Split(ormTag, ",")
	tagMap := make(map[string]string, len(pairs))

	for _, pair := range pairs {
		content := strings.Split(pair, "=")
		if len(content) != 2 {
			return nil, errs.InvalidTagContentErr(pair)
		}

		key := strings.Trim(content[0], " ")
		if key == "" {
			return nil, errs.EmptyTagKeyErr(pair)
		}

		val := strings.Trim(content[1], " ")
		if val == "" {
			return nil, errs.EmptyTagValErr(pair)
		}

		tagMap[key] = val
	}

	return tagMap, nil
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
