package workflow

import "github.com/creatortsv/workflow-go/pkg/state"

type builder[T any, E comparable] struct {
	sm state.Manager[T, E]
	ts map[string]Transition[T, E]
}

func NewBuilder[T any, E comparable](m state.Manager[T, E]) *builder[T, E] {
	return &builder[T, E]{
		sm: m,
		ts: map[string]Transition[T, E]{},
	}
}

func (b *builder[T, E]) WithTransition(
	name string,
	dist E,
	from []E,
	guards ...GuardFunc[T],
) *builder[T, E] {
	b.ts[name] = &transition[T, E]{
		guards: guards,
		name:   name,
		from:   from,
		dist:   dist,
	}

	return b
}

func (b *builder[T, E]) Build() *workflow[T, E] {
	return &workflow[T, E]{
		sm: b.sm,
		ts: b.ts,
	}
}
