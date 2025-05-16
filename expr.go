package easyorm

// Expr marker interface, representing an expression.
// All the elements after "WHERE" are expression.
type Expr interface {
	expr()
}
