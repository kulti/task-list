package models_test

import (
	"errors"
	"testing"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/stretchr/testify/require"
)

type transition map[models.SwitchTaskStateEvent]models.TaskState

func TestNextTaskState(t *testing.T) {
	transitions := map[models.TaskState]transition{
		models.TaskStateSimple: {
			models.TodoTaskEvent:     models.TaskStateTodo,
			models.DoneTaskEvent:     models.TaskStateCompleted,
			models.UndoneTaskEvent:   models.TaskStateSimple,
			models.CancelTaskEvent:   models.TaskStateCanceled,
			models.PostponeTaskEvent: models.TaskStateSimple,
		},
		models.TaskStateTodo: {
			models.TodoTaskEvent:     models.TaskStateTodo,
			models.DoneTaskEvent:     models.TaskStateCompleted,
			models.UndoneTaskEvent:   models.TaskStateTodo,
			models.CancelTaskEvent:   models.TaskStateCanceled,
			models.PostponeTaskEvent: models.TaskStateSimple,
		},
		models.TaskStateCompleted: {
			models.DoneTaskEvent:   models.TaskStateCompleted,
			models.UndoneTaskEvent: models.TaskStateSimple,
		},
		models.TaskStateCanceled: {
			models.CancelTaskEvent: models.TaskStateCanceled,
			models.ToWorkTaskEvent: models.TaskStateSimple,
		},
	}

	for state, tr := range transitions {
		for ev, expNextState := range tr {
			state := state
			ev := ev
			expNextState := expNextState
			t.Run(string(state)+" -> "+ev.String(), func(t *testing.T) {
				t.Parallel()
				nextState, err := state.NextState(ev)
				require.NoError(t, err)
				require.Equal(t, expNextState, nextState)
			})
		}
	}
}

func TestNextTaskStateInconcistency(t *testing.T) {
	unknownState := models.TaskState(-1)
	unknownEvent := models.SwitchTaskStateEvent(-1)

	transitions := map[models.TaskState]map[models.SwitchTaskStateEvent]struct{}{
		unknownState: {
			models.TodoTaskEvent:   struct{}{},
			models.DoneTaskEvent:   struct{}{},
			models.UndoneTaskEvent: struct{}{},
			models.CancelTaskEvent: struct{}{},
			models.ToWorkTaskEvent: struct{}{},
			unknownEvent:           struct{}{},
		},
		models.TaskStateSimple: {
			models.ToWorkTaskEvent: struct{}{},
			unknownEvent:           struct{}{},
		},
		models.TaskStateTodo: {
			models.ToWorkTaskEvent: struct{}{},
			unknownEvent:           struct{}{},
		},
		models.TaskStateCompleted: {
			models.TodoTaskEvent:     struct{}{},
			models.CancelTaskEvent:   struct{}{},
			models.PostponeTaskEvent: struct{}{},
			models.ToWorkTaskEvent:   struct{}{},
			unknownEvent:             struct{}{},
		},
		models.TaskStateCanceled: {
			models.TodoTaskEvent:     struct{}{},
			models.DoneTaskEvent:     struct{}{},
			models.UndoneTaskEvent:   struct{}{},
			models.PostponeTaskEvent: struct{}{},
			unknownEvent:             struct{}{},
		},
	}

	for state, tr := range transitions {
		for ev := range tr {
			state := state
			ev := ev
			t.Run(string(state)+" -> "+ev.String(), func(t *testing.T) {
				t.Parallel()
				_, err := state.NextState(ev)
				require.Error(t, err)
				require.True(t, errors.As(err, &models.StateInconsistencyErr{}))
			})
		}
	}
}
