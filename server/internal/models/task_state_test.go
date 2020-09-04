package models_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kulti/task-list/server/internal/models"
)

func TestValidateStateSwitch(t *testing.T) {
	transitions := map[models.TaskState][]models.SwitchTaskStateEvent{
		models.TaskStateSimple: {
			models.TodoTaskEvent,
			models.DoneTaskEvent,
			models.CancelTaskEvent,
			models.PostponeTaskEvent,
		},
		models.TaskStateTodo: {
			models.TodoTaskEvent,
			models.DoneTaskEvent,
			models.CancelTaskEvent,
			models.PostponeTaskEvent,
		},
		models.TaskStateCompleted: {
			models.DoneTaskEvent,
			models.ToWorkTaskEvent,
		},
		models.TaskStateCanceled: {
			models.CancelTaskEvent,
			models.ToWorkTaskEvent,
		},
	}

	for state, tr := range transitions {
		for _, ev := range tr {
			state := state
			ev := ev
			t.Run(string(state)+" -> "+ev.String(), func(t *testing.T) {
				t.Parallel()
				require.NoError(t, state.ValidateStateSwitch(ev))
			})
		}
	}
}

func TestValidateStateSwitchInconcistency(t *testing.T) {
	unknownState := models.TaskState(-1)
	unknownEvent := models.SwitchTaskStateEvent(-1)

	transitions := map[models.TaskState][]models.SwitchTaskStateEvent{
		unknownState: {
			models.TodoTaskEvent,
			models.DoneTaskEvent,
			models.CancelTaskEvent,
			models.ToWorkTaskEvent,
			unknownEvent,
		},
		models.TaskStateSimple: {
			models.ToWorkTaskEvent,
			unknownEvent,
		},
		models.TaskStateTodo: {
			models.ToWorkTaskEvent,
			unknownEvent,
		},
		models.TaskStateCompleted: {
			models.TodoTaskEvent,
			models.CancelTaskEvent,
			models.PostponeTaskEvent,
			unknownEvent,
		},
		models.TaskStateCanceled: {
			models.TodoTaskEvent,
			models.DoneTaskEvent,
			models.PostponeTaskEvent,
			unknownEvent,
		},
	}

	for state, tr := range transitions {
		for _, ev := range tr {
			state := state
			ev := ev
			t.Run(string(state)+" -> "+ev.String(), func(t *testing.T) {
				t.Parallel()
				err := state.ValidateStateSwitch(ev)
				require.Error(t, err)
				require.True(t, errors.As(err, &models.StateInconsistencyErr{}))
			})
		}
	}
}
