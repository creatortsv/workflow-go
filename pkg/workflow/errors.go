package workflow

import "errors"

var (
	ErrUnknownTransition   = errors.New("unknown transition name")
	ErrForbiddenTransition = errors.New("transition is forbidden")
)
