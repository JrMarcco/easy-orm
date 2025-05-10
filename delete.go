package easyorm

type Deleter[T any] struct {
	builder

	session session
}
