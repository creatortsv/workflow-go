package workflow

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	testManager struct{}
	testSubject struct {
		state string
	}
)

func (m *testManager) State(_ context.Context, subject *testSubject) (string, error) {
	return subject.state, nil
}

func (m *testManager) Place(_ context.Context, subject *testSubject, state string) error {
	subject.state = state

	return nil
}

func TestWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("New: ok", func(t *testing.T) {
		t.Parallel()

		m := &testManager{}
		w := NewBuilder(m).WithTransition("start", "started", []string{"init"}).Build()

		require.Equal(t, m, w.sm)

		i, ok := w.ts["start"]

		require.True(t, ok)
		require.Equal(t, "start", i.String())
		require.Equal(t, "started", i.Dist())
		require.Equal(t, []string{"init"}, i.From())
	})

	t.Run("Apply: unknown transition", func(t *testing.T) {
		t.Parallel()

		err := NewBuilder(&testManager{}).Build().Apply(context.Background(), &testSubject{}, "start")

		require.ErrorIs(t, err, ErrUnknownTransition)
	})

	t.Run("Apply: forbidden transition by state", func(t *testing.T) {
		t.Parallel()

		w := NewBuilder(&testManager{}).WithTransition("start", "started", []string{"init"}).Build()

		err := w.Apply(context.Background(), &testSubject{state: "started"}, "start")

		require.ErrorIs(t, err, ErrForbiddenTransition)
	})

	t.Run("Apply: forbidden transition by guards", func(t *testing.T) {
		t.Parallel()

		w := NewBuilder(&testManager{}).WithTransition(
			"start",
			"started",
			[]string{"init"},
			func(_ context.Context, _ *testSubject) bool { return true },
			func(_ context.Context, _ *testSubject) bool { return false },
		).Build()

		err := w.Apply(context.Background(), &testSubject{state: "init"}, "start")

		require.ErrorIs(t, err, ErrForbiddenTransition)
	})

	t.Run("Apply: ok", func(t *testing.T) {
		t.Parallel()

		subject := &testSubject{state: "init"}

		w := NewBuilder(&testManager{}).WithTransition(
			"start",
			"started",
			[]string{"init"},
			func(_ context.Context, _ *testSubject) bool { return true },
		).Build()

		err := w.Apply(context.Background(), subject, "start")

		require.NoError(t, err)
		require.Equal(t, "started", subject.state)
	})

	t.Run("AllowedTransitions: ok", func(t *testing.T) {
		t.Parallel()

		subject := &testSubject{state: "init"}

		w := NewBuilder(&testManager{}).WithTransition(
			"start",
			"started",
			[]string{"init"},
			func(_ context.Context, _ *testSubject) bool { return true },
		).Build()

		c, err := w.AllowedTransitions(context.Background(), subject)

		require.NoError(t, err)
		require.True(t, slices.Contains(c, "start"))
	})
}
