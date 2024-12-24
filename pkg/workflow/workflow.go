package workflow

import (
	"context"
	"fmt"
	"slices"

	"github.com/creatortsv/workflow-go/pkg/state"
)

type workflow[T any, E comparable] struct {
	sm state.Manager[T, E]
	ts map[string]ReadOnlyTransition[T, E]
}

type withTransition[T any, E comparable] func() *transition[T, E]

func New[T any, E comparable](s state.Manager[T, E]) *workflow[T, E] {
	return &workflow[T, E]{
		sm: s,
		ts: map[string]ReadOnlyTransition[T, E]{},
	}
}

func (w *workflow[T, E]) WithTransition(
	name string,
	dist E,
	from []E,
	guards ...GuardFunc[T],
) *workflow[T, E] {
	t := &transition[T, E]{
		guards: guards,
		name:   name,
		from:   from,
		dist:   dist,
	}

	w.ts[t.name] = t

	return w
}

func (w *workflow[T, E]) Apply(ctx context.Context, subject T, transition string) error {
	if _, ok := w.ts[transition]; !ok {
		return fmt.Errorf("apply transition [%s]: %w", transition, ErrUnknownTransition)
	}

	t, err := w.AllowedTransitions(ctx, subject)
	if err != nil {
		return err
	}

	if t, ok := t[transition]; ok {
		if err := w.sm.Place(ctx, subject, t.Dist()); err != nil {
			return fmt.Errorf("apply transition [%s]: %w", transition, err)
		}

		return nil
	}

	return fmt.Errorf("apply transition [%s]: %w", transition, ErrForbiddenTransition)
}

func (w *workflow[T, E]) AllowedTransitions(ctx context.Context, subject T) (map[string]ReadOnlyTransition[T, E], error) {
	s, err := w.sm.State(ctx, subject)
	if err != nil {
		return nil, fmt.Errorf("getting current state: %w", err)
	}

	transits := make(map[string]ReadOnlyTransition[T, E])
	for _, t := range w.ts {
		if !slices.Contains(t.From(), s) {
			continue
		}

		if !t.IsAllowed(ctx, subject) {
			continue
		}

		transits[t.String()] = t
	}

	return transits, nil
}
