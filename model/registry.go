package model

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/JrMarcco/easy-orm/internal/errs"
)

const (
	tagName    = "orm"
	tagNameCol = "column"
)

var _ Registry = (*modelRegistry)(nil)

type modelRegistry struct {
	sync.RWMutex
	models map[reflect.Type]*Model
}

func (r *modelRegistry) GetModel(entity any) (*Model, error) {
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

	m, err := r.parseModel(typ)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m

	return m, nil
}

// RegisterModel register a model.
func (r *modelRegistry) RegisterModel(entity any, opts ...Opt) (*Model, error) {
	typ := reflect.TypeOf(entity)
	m, err := r.parseModel(typ)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}

	r.models[typ] = m
	return m, nil
}

func NewRegistry() Registry {
	return &modelRegistry{
		models: make(map[reflect.Type]*Model, 32),
	}
}

// parseModel parse the model from the given entity.
//
// an entity must be a struct or a pointer to struct.
func (r *modelRegistry) parseModel(typ reflect.Type) (*Model, error) {
	elemTyp := typ
	if elemTyp.Kind() != reflect.Struct {
		if elemTyp.Kind() != reflect.Pointer {
			return nil, errs.ErrInvalidModelType
		}

		elemTyp = elemTyp.Elem()

		// only support a pointer to struct
		if elemTyp.Kind() != reflect.Struct {
			return nil, errs.ErrInvalidModelType
		}
	}

	numField := elemTyp.NumField()
	fields := make(map[string]*Field, numField)
	columns := make(map[string]*Field, numField)

	for i := 0; i < numField; i++ {
		structField := elemTyp.Field(i)

		// parse tag
		tagMap, err := r.parseTag(structField.Tag)
		if err != nil {
			return nil, err
		}

		colName, ok := tagMap[tagNameCol]
		if !ok {
			colName = camelToUnderline(structField.Name)
		}

		field := &Field{
			Typ:        structField.Type,
			FiledName:  structField.Name,
			ColumnName: colName,
			Offset:     structField.Offset,
		}
		fields[structField.Name] = field
		columns[colName] = field
	}

	return &Model{
		TableName: camelToUnderline(elemTyp.Name()),
		Fields:    fields,
		Columns:   columns,
	}, nil
}

// parseTag parse the orm tag from the given tag.
//
// like this:
//
//	type User struct {
//		Name string `orm:"column=user_name"`
//	}
//
// the tag content is "column=user_id,column=user_name"
// the tag name is "column"
// the tag value is "user_name"
func (r *modelRegistry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup(tagName)
	if !ok {
		return map[string]string{}, nil
	}

	pairs := strings.Split(ormTag, ",")
	tagMap := make(map[string]string, len(pairs))

	for _, pair := range pairs {
		content := strings.Split(pair, "=")
		if len(content) != 2 {
			return nil, errs.ErrInvalidTag(pair)
		}

		key := strings.Trim(content[0], " ")
		if key == "" {
			return nil, errs.ErrInvalidTag(pair)
		}

		val := strings.Trim(content[1], " ")
		if val == "" {
			return nil, errs.ErrInvalidTag(pair)
		}

		tagMap[key] = val
	}

	return tagMap, nil
}

func camelToUnderline(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(s[i-1])
			// the previous character is lower case
			// or
			// the next character is lower case
			if (prev >= 'a' && prev <= 'z') || (i+1 < len(s) && s[i+1] >= 'a' && s[i+1] <= 'z') {
				result = append(result, '_')
			}
		}
		if r == '-' || r == ' ' || r == '.' {
			result = append(result, '_')
		} else {
			result = append(result, r)
		}
	}
	// turn to the lower case
	out := strings.ToLower(string(result))

	// remove redundant underscores
	out = strings.Trim(out, "_")
	out = regexp.MustCompile(`_+`).ReplaceAllString(out, "_")
	return out
}
