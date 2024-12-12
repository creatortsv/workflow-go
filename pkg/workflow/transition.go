package workflow

import (
	"context"
)

type transition[T any] struct {
	name   string
	guards []GuardFunc[T]
	from   []string
	dist   string
}

func (t *transition[T]) IsAllowed(ctx context.Context, subject T) bool {
	for _, g := range t.guards {
		if g.IsAllowed(ctx, subject) {
			continue
		}

		return false
	}

	return true
}

func (t *transition[T]) String() string {
	return t.name
}

func (t *transition[T]) Dist() string {
	return t.dist
}

func (t *transition[T]) From() []string {
	return t.from
}
