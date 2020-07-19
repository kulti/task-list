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
			models.TodoTaskEvent:   models.TaskStateTodo,
			models.DoneTaskEvent:   models.TaskStateCompleted,
			models.UndoneTaskEvent: models.TaskStateSimple,
			models.CancelTaskEvent: models.TaskStateCanceled,
		},
		models.TaskStateTodo: {
			models.TodoTaskEvent:   models.TaskStateTodo,
			models.DoneTaskEvent:   models.TaskStateCompleted,
			models.UndoneTaskEvent: models.TaskStateTodo,
			models.CancelTaskEvent: models.TaskStateCanceled,
		},
		models.TaskStateCompleted: {
			models.DoneTaskEvent:   models.TaskStateCompleted,
			models.UndoneTaskEvent: models.TaskStateSimple,
		},
		models.TaskStateCanceled: {
			models.CancelTaskEvent: models.TaskStateCanceled,
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
			unknownEvent:           struct{}{},
		},
		models.TaskStateSimple: {
			unknownEvent: struct{}{},
		},
		models.TaskStateTodo: {
			unknownEvent: struct{}{},
		},
		models.TaskStateCompleted: {
			models.TodoTaskEvent:   struct{}{},
			models.CancelTaskEvent: struct{}{},
			unknownEvent:           struct{}{},
		},
		models.TaskStateCanceled: {
			models.TodoTaskEvent:   struct{}{},
			models.DoneTaskEvent:   struct{}{},
			models.UndoneTaskEvent: struct{}{},
			unknownEvent:           struct{}{},
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
