package storages

import (
	"context"

	"github.com/kulti/task-list/server/internal/models"
)

// TaskStore is an interface to task storage.
type TaskStore interface {
	NewSprint(ctx context.Context, opts models.SprintOpts) error
	CreateTask(ctx context.Context, task models.Task, listID string) (string, error)
	DeleteTaskFromList(ctx context.Context, taskID, listID string) error
	ListTasks(ctx context.Context, listID string) (models.TaskList, error)
	UpdateTask(ctx context.Context, taskID string, points models.UpdateOptions) error
	TodoTask(ctx context.Context, taskID string) error
	DoneTask(ctx context.Context, taskID string) error
	CancelTask(ctx context.Context, taskID string) error
	UndoneTask(ctx context.Context, taskID string) error
	PostponeTask(ctx context.Context, taskID string) error
	GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error)
}
