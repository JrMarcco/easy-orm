package orm

// Assignable 标记接口，表示可赋值
type Assignable interface {
	assign()
}

type Assignment struct {
	fdName string
	val    any
}

func (a Assignment) assign() {}

func Assign(fdName string, val any) Assignment {
	return Assignment{
		fdName: fdName,
		val:    val,
	}
}
