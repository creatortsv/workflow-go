package workflow

import (
	"context"
	"fmt"
)

type (
	Workflow[T any, E comparable] interface {
		Apply(ctx context.Context, subject T, transition string) error
		State(ctx context.Context, subject T) (E, error)
		AllowedTransitions(ctx context.Context, subject T) ([]string, error)
	}

	Transition[T any, E comparable] interface {
		Guard[T]
		fmt.Stringer
		Dist() E
		From() []E
	}

	GuardFunc[T any] func(context.Context, T) bool
	Guard[T any]     interface {
		IsAllowed(ctx context.Context, subject T) bool
	}
)

func (f GuardFunc[T]) IsAllowed(ctx context.Context, subject T) bool {
	return f(ctx, subject)
}
