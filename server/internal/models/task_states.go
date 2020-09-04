package models

// SwitchTaskStateEvent represents events that can change task state.
type SwitchTaskStateEvent int

// Switch task states events.
const (
	DoneTaskEvent SwitchTaskStateEvent = iota
	TodoTaskEvent
	CancelTaskEvent
	PostponeTaskEvent
	ToWorkTaskEvent
)

func (e SwitchTaskStateEvent) String() string {
	switch e {
	case DoneTaskEvent:
		return "[done]"
	case TodoTaskEvent:
		return "[todo]"
	case CancelTaskEvent:
		return "[cancel]"
	case PostponeTaskEvent:
		return "[postpone]"
	case ToWorkTaskEvent:
		return "[towork]"
	default:
		return "[unknown]"
	}
}

// ValidateStateSwitch check if task state can be switched by event.
func (s TaskState) ValidateStateSwitch(ev SwitchTaskStateEvent) error {
	switch s {
	case TaskStateSimple, TaskStateTodo:
		switch ev {
		case DoneTaskEvent, TodoTaskEvent, CancelTaskEvent, PostponeTaskEvent:
			return nil
		case ToWorkTaskEvent:
			return NewStateInconsistencyErr(s, ev)
		}
	case TaskStateCompleted:
		switch ev {
		case DoneTaskEvent, ToWorkTaskEvent:
			return nil
		case CancelTaskEvent, PostponeTaskEvent, TodoTaskEvent:
			return NewStateInconsistencyErr(s, ev)
		}
	case TaskStateCanceled:
		switch ev {
		case ToWorkTaskEvent, CancelTaskEvent:
			return nil
		case DoneTaskEvent, PostponeTaskEvent, TodoTaskEvent:
			return NewStateInconsistencyErr(s, ev)
		}
	}
	return NewStateInconsistencyErr(s, ev)
}
