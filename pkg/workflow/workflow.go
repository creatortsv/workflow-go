package workflow

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/creatortsv/workflow-go/pkg/state"
)

var _ Workflow[any, string] = &workflow[any, string]{}

type workflow[T any, E comparable] struct {
	sm state.Manager[T, E]
	ts map[string]Transition[T, E]
}

func (w *workflow[T, E]) Apply(ctx context.Context, subject T, transition string) error {
	if _, ok := w.ts[transition]; !ok {
		return fmt.Errorf("apply transition [%s]: %w", transition, ErrUnknownTransition)
	}

	t, err := w.allowedTransitions(ctx, subject)
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

func (w *workflow[T, E]) State(ctx context.Context, subject T) (E, error) {
	return w.sm.State(ctx, subject)
}

func (w *workflow[T, E]) AllowedTransitions(ctx context.Context, subject T) ([]string, error) {
	c, err := w.allowedTransitions(ctx, subject)
	if err != nil {
		return nil, err
	}

	return slices.Collect(maps.Keys(c)), nil
}

func (w *workflow[T, E]) allowedTransitions(ctx context.Context, subject T) (map[string]Transition[T, E], error) {
	s, err := w.sm.State(ctx, subject)
	if err != nil {
		return nil, fmt.Errorf("getting current state: %w", err)
	}

	transits := make(map[string]Transition[T, E])
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
