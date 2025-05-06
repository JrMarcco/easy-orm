package easyorm

// Expression marker interface, representing an expression.
// All the elements after "WHERE" are expression.
type Expression interface {
	expr()
}
