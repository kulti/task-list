package storages

import (
	"time"

	"github.com/kulti/task-list/server/internal/models"
)

// SprintOpts represents new sprint options.
type SprintOpts struct {
	Title string
	Begin time.Time
	End   time.Time
}

// TaskList represents a task list.
type TaskList struct {
	Title string
	Tasks []Task
}

// Task represents a task.
type Task struct {
	ID     int64
	Text   string
	State  models.TaskState
	Points int32
	Burnt  int32
}

// UpdateTaskFn should returns an updated task or error.
type UpdateTaskFn func(Task) (Task, error)

// PostponeTaskFn should returns a postponend task and probably updated task or error.
// Postpone operation can split task into two. The first part should be postponed,
// the second should be updated.
type PostponeTaskFn func(Task) (Task, Task, error)
