package orm

// Expression 表达式
// 可以这样理解，跟在 WHERE 后的所有元素的都是表达式。
// 其最终会构建成一个表达式的二叉树。
type Expression interface {
	expr()
}

// RawExpr 原生表达式
type RawExpr struct {
	raw  string
	args []any
}

var _ Expression = new(RawExpr)
var _ selectable = new(RawExpr)

func (r RawExpr) expr()       {}
func (r RawExpr) selectable() {}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		left: r,
	}
}

func Raw(expr string, args ...any) RawExpr {
	return RawExpr{
		raw:  expr,
		args: args,
	}
}
