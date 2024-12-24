package workflow

import (
	"context"
)

type transition[T any, E comparable] struct {
	name   string
	guards []GuardFunc[T]
	from   []E
	dist   E
}

func (t *transition[T, E]) IsAllowed(ctx context.Context, subject T) bool {
	for _, g := range t.guards {
		if g.IsAllowed(ctx, subject) {
			continue
		}

		return false
	}

	return true
}

func (t *transition[T, E]) String() string {
	return t.name
}

func (t *transition[T, E]) Dist() E {
	return t.dist
}

func (t *transition[T, E]) From() []E {
	return t.from
}
