package easyorm

type Assignable interface {
	assign()
}

var _ Assignable = (*Assignment)(nil)

type Assignment struct {
	filedName string
	value     any
}

func (a Assignment) assign() {}

func Assign(fieldName string, value any) Assignment {
	return Assignment{
		filedName: fieldName,
		value:     value,
	}
}
