package models

// SprintOpts represents new sprint options.
type SprintOpts struct {
	Title string `json:"title"`
	Begin string `json:"begin"`
	End   string `json:"end"`
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

// UpdateOptions represents an update options.
type UpdateOptions struct {
	Text   string `json:"text"`
	Points int32  `json:"points"`
	Burnt  int32  `json:"burnt"`
}

// SprintTemplate represents a sprint template.
type SprintTemplate struct {
	Tasks []TaskTemplate `json:"tasks"`
}

// TaskTemplate represents a task template.
type TaskTemplate struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Points int32  `json:"points"`
}
