package models

// SwitchTaskStateEvent represents events that can change task state.
type SwitchTaskStateEvent int

// Switch task states events.
const (
	DoneTaskEvent SwitchTaskStateEvent = iota
	UndoneTaskEvent
	TodoTaskEvent
	CancelTaskEvent
	PostponeTaskEvent
	ToWorkTaskEvent
)

func (e SwitchTaskStateEvent) String() string {
	switch e {
	case DoneTaskEvent:
		return "[done]"
	case UndoneTaskEvent:
		return "[undone]"
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

// NextState returns next task state caused by event.
func (s TaskState) NextState(ev SwitchTaskStateEvent) (TaskState, error) {
	switch s {
	case TaskStateSimple, TaskStateTodo:
		switch ev {
		case DoneTaskEvent:
			return TaskStateCompleted, nil
		case UndoneTaskEvent:
			return s, nil
		case TodoTaskEvent:
			return TaskStateTodo, nil
		case CancelTaskEvent:
			return TaskStateCanceled, nil
		case PostponeTaskEvent:
			return TaskStateSimple, nil
		case ToWorkTaskEvent:
			return TaskStateSimple, NewStateInconsistencyErr(s, ev)
		}
	case TaskStateCompleted:
		switch ev {
		case UndoneTaskEvent:
			return TaskStateSimple, nil
		case DoneTaskEvent:
			return TaskStateCompleted, nil
		case CancelTaskEvent, PostponeTaskEvent, TodoTaskEvent, ToWorkTaskEvent:
			return TaskStateSimple, NewStateInconsistencyErr(s, ev)
		}
	case TaskStateCanceled:
		switch ev {
		case CancelTaskEvent:
			return TaskStateCanceled, nil
		case ToWorkTaskEvent:
			return TaskStateSimple, nil
		case DoneTaskEvent, PostponeTaskEvent, TodoTaskEvent, UndoneTaskEvent:
			return TaskStateSimple, NewStateInconsistencyErr(s, ev)
		}
	}
	return TaskStateSimple, NewStateInconsistencyErr(s, ev)
}
