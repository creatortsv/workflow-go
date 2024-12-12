package state

import (
	"context"
)

type Manager[T any] interface {
	// Returns the current state of subject
	State(ctx context.Context, subject T) (string, error)

	// Sets the given state
	Place(ctx context.Context, subject T, state string) error
}
