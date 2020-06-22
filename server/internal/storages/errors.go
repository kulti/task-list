package storages

import "fmt"

type StateInconsistencyErr struct {
	fromState, toState string
}

func NewStateInconsistencyErr(fromState, toState string) StateInconsistencyErr {
	return StateInconsistencyErr{fromState: fromState, toState: toState}
}

func (e StateInconsistencyErr) Error() string {
	return fmt.Sprintf("cannot switch state from %s to %s", e.fromState, e.toState)
}
