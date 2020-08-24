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
