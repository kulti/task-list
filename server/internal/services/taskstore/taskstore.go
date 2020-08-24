package taskstore

import (
	"context"
	"strconv"

	"github.com/kulti/task-list/server/internal/models"
)

type dbStore interface {
	DeleteTask(ctx context.Context, taskID int64) error
	UpdateTask(ctx context.Context, taskID int64, points models.UpdateOptions) error
	TodoTask(ctx context.Context, taskID int64) error
	DoneTask(ctx context.Context, taskID int64) error
	CancelTask(ctx context.Context, taskID int64) error
	BackTaskToWork(ctx context.Context, taskID int64) error
	UndoneTask(ctx context.Context, taskID int64) error
	PostponeTask(ctx context.Context, taskID int64) error
}

// TaskStore provides methods to manage tasks.
type TaskStore struct {
	dbStore dbStore
}

// New creates a new instance of TaskStore.
func New(dbStore dbStore) *TaskStore {
	return &TaskStore{
		dbStore: dbStore,
	}
}

// DeleteTask deletes task.
func (s *TaskStore) DeleteTask(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.DeleteTask(ctx, taskID)
	})
}

// UpdateTask updates task text and points.
func (s *TaskStore) UpdateTask(ctx context.Context, taskID string, opts models.UpdateOptions) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.UpdateTask(ctx, taskID, opts)
	})
}

// TodoTask changes task state to todo.
func (s *TaskStore) TodoTask(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.TodoTask(ctx, taskID)
	})
}

// DoneTask changes task burnt points to be equal all points.
func (s *TaskStore) DoneTask(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.DoneTask(ctx, taskID)
	})
}

// CancelTask changes task state to canceled.
func (s *TaskStore) CancelTask(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.CancelTask(ctx, taskID)
	})
}

// BackTaskToWork changes task state to "" from canceled.
func (s *TaskStore) BackTaskToWork(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.BackTaskToWork(ctx, taskID)
	})
}

// UndoneTask chage task state to "". Crappy method - will be removed soon.
func (s *TaskStore) UndoneTask(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.UndoneTask(ctx, taskID)
	})
}

// PostponeTask postpones task to the next sprint.
// It looks like to cancel and move task to the next sprint template.
func (s *TaskStore) PostponeTask(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.PostponeTask(ctx, taskID)
	})
}

func (s *TaskStore) doWithTaskIDConvert(taskID string, fn func(taskID int64) error) error {
	id, err := strconv.ParseInt(taskID, 16, 64)
	if err != nil {
		return err
	}

	return fn(id)
}
