package models

// SprintOpts represents new sprint options.
type SprintOpts struct {
	Title string `json:"title"`
}

// TaskList represents a task list.
type TaskList struct {
	Title string `json:"title"`
	Tasks []Task `json:"tasks"`
}

// Task represents a task.
type Task struct {
	ID     string    `json:"id"`
	Text   string    `json:"text"`
	State  TaskState `json:"state"`
	Points int32     `json:"points"`
	Burnt  int32     `json:"burnt,omitempty"`
}

// TaskState reprensts a task state.
type TaskState string

// TaskState constants.
const (
	TaskStateTodo      TaskState = "todo"
	TaskStateCompleted TaskState = "completed"
	TaskStateCanceled  TaskState = "canceled"
)

// UpdateOptions represents a update options.
type UpdateOptions struct {
	Text   string `json:"text"`
	Points int32  `json:"points"`
	Burnt  int32  `json:"burnt"`
}
