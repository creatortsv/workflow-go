package workflow

import (
	"context"
	"fmt"
	"slices"

	"github.com/creatortsv/workflow-go/pkg/state"
)

type workflow[T any] struct {
	sm state.Manager[T]
	ts map[string]ReadOnlyTransition[T]
}

type withTransition[T any] func() *transition[T]

func New[T any](s state.Manager[T]) *workflow[T] {
	return &workflow[T]{
		sm: s,
		ts: make(map[string]ReadOnlyTransition[T]),
	}
}

func (w *workflow[T]) WithTransition(
	name string,
	dist string,
	from []string,
	guards ...GuardFunc[T],
) *workflow[T] {
	t := &transition[T]{
		guards: guards,
		name:   name,
		from:   from,
		dist:   dist,
	}

	w.ts[t.name] = t

	return w
}

func (w *workflow[T]) Apply(ctx context.Context, subject T, transition string) error {
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

func (w *workflow[T]) AllowedTransitions(ctx context.Context, subject T) (map[string]ReadOnlyTransition[T], error) {
	s, err := w.sm.State(ctx, subject)
	if err != nil {
		return nil, fmt.Errorf("getting current state: %w", err)
	}

	transits := make(map[string]ReadOnlyTransition[T])
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
