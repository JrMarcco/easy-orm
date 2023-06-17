package orm

// Expression 表达式
// 可以这样理解，跟在 WHERE 后的所有元素的都是表达式。
// 其最终会构建成一个表达式的二叉树。
type Expression interface {
	expr()
}
