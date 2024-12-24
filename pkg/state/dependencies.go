package state

import (
	"context"
)

type Manager[T any, E comparable] interface {
	// Returns the current state of subject
	State(ctx context.Context, subject T) (E, error)

	// Sets the given state
	Place(ctx context.Context, subject T, state E) error
}
