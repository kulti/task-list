package storages

import (
	"context"

	"github.com/kulti/task-list/internal/models"
)

// TaskStore is an interface to task storage.
type TaskStore interface {
	NewSprint(ctx context.Context, title string) error
	CreateTask(ctx context.Context, task models.Task, listIDs []string) (string, error)
	TakeTaskToList(ctx context.Context, taskID, listIDs string) error
	DeleteTaskFromList(ctx context.Context, taskID, listID string) error
	ListTasks(ctx context.Context, listID string) (models.TaskList, error)
	UpdateTask(ctx context.Context, taskID string, points models.UpdateOptions) error
	DoneTask(ctx context.Context, taskID string) error
	CancelTask(ctx context.Context, taskID string) error
}
