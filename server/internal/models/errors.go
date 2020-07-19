package models

import (
	"fmt"
)

// StateInconsistencyErr is an error when state cannot be switched by event.
type StateInconsistencyErr struct {
	state TaskState
	event SwitchTaskStateEvent
}

// NewStateInconsistencyErr creates a new instance of StateInconsistencyErr.
func NewStateInconsistencyErr(state TaskState, event SwitchTaskStateEvent) StateInconsistencyErr {
	return StateInconsistencyErr{state: state, event: event}
}

func (e StateInconsistencyErr) Error() string {
	return fmt.Sprintf("cannot switch state from %s by event %s", e.state, e.event)
}
