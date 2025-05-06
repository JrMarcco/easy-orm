package model

import (
	"github.com/JrMarcco/easy-orm/internal/errs"
	"reflect"
	"regexp"
	"strings"
	"sync"
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

	m, err := parseModel(entity)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (r *modelRegistry) RegisterModel(entity any, opts ...Opt) (*Model, error) {
	m, err := parseModel(entity)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(m)
	}

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
func parseModel(entity any) (*Model, error) {
	typ := reflect.TypeOf(entity)

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

	for i := 0; i < numField; i++ {
		field := elemTyp.Field(i)

		fields[field.Name] = &Field{
			FiledName:  field.Name,
			ColumnName: camelToUnderline(field.Name),
		}

	}
	return &Model{
		TableName: camelToUnderline(elemTyp.Name()),
		Fields:    fields,
	}, nil
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
