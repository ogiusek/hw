package hw

import "net/http"

type Executable[T any] func(args T) any

func Run[T any](method Executable[T]) Wrapper[T] {
	return func(args T, r *http.Request) any { return method(args) }
}
