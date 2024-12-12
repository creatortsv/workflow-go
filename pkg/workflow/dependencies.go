package workflow

import (
	"context"
	"fmt"
)

type ReadOnlyTransition[T any] interface {
	fmt.Stringer
	Guard[T]
	Dist() string
	From() []string
}

type (
	GuardFunc[T any] func(context.Context, T) bool
	Guard[T any]     interface {
		IsAllowed(ctx context.Context, subject T) bool
	}
)

func (f GuardFunc[T]) IsAllowed(ctx context.Context, subject T) bool {
	return f(ctx, subject)
}
