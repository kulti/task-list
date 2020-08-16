package storages

import (
	"context"

	"github.com/kulti/task-list/server/internal/models"
)

// TaskStore is an interface to task storage.
type TaskStore interface {
	NewSprint(ctx context.Context, opts SprintOpts) error
	CreateTask(ctx context.Context, task models.Task, sprintID string) (string, error)
	DeleteTask(ctx context.Context, taskID string) error
	ListTasks(ctx context.Context, sprintID string) (models.TaskList, error)
	UpdateTask(ctx context.Context, taskID string, points models.UpdateOptions) error
	TodoTask(ctx context.Context, taskID string) error
	DoneTask(ctx context.Context, taskID string) error
	CancelTask(ctx context.Context, taskID string) error
	BackTaskToWork(ctx context.Context, taskID string) error
	UndoneTask(ctx context.Context, taskID string) error
	PostponeTask(ctx context.Context, taskID string) error
	GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error)
}
