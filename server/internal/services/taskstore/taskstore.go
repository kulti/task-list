package taskstore

import (
	"context"
	"strconv"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

type dbStore interface {
	UpdateTask(ctx context.Context, taskID int64, fn storages.UpdateTaskFn) error
	PostponeTask(ctx context.Context, taskID int64, fn storages.PostponeTaskFn) error
	DeleteTask(ctx context.Context, taskID int64) error
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
	return s.updateWithTaskIDConvert(ctx, taskID, func(task storages.Task) (storages.Task, error) {
		newState := task.State
		if opts.Points == opts.Burnt {
			newState = models.TaskStateCompleted
			if err := task.State.ValidateStateSwitch(models.DoneTaskEvent); err != nil {
				return storages.Task{}, err
			}
		} else if task.State == models.TaskStateCompleted {
			newState = models.TaskStateSimple
			if err := task.State.ValidateStateSwitch(models.ToWorkTaskEvent); err != nil {
				return storages.Task{}, err
			}
		}

		newTask := storages.Task{
			ID:     task.ID,
			State:  newState,
			Text:   opts.Text,
			Points: opts.Points,
			Burnt:  opts.Burnt,
		}
		return newTask, nil
	})
}

// TodoTask changes task state to todo.
func (s *TaskStore) TodoTask(ctx context.Context, taskID string) error {
	return s.updateWithTaskIDConvert(ctx, taskID, func(task storages.Task) (storages.Task, error) {
		if err := task.State.ValidateStateSwitch(models.TodoTaskEvent); err != nil {
			return storages.Task{}, err
		}

		newTask := task
		newTask.State = models.TaskStateTodo
		return newTask, nil
	})
}

// DoneTask changes task burnt points to be equal all points.
func (s *TaskStore) DoneTask(ctx context.Context, taskID string) error {
	return s.updateWithTaskIDConvert(ctx, taskID, func(task storages.Task) (storages.Task, error) {
		if err := task.State.ValidateStateSwitch(models.DoneTaskEvent); err != nil {
			return storages.Task{}, err
		}

		newTask := task
		newTask.State = models.TaskStateCompleted
		newTask.Burnt = newTask.Points
		return newTask, nil
	})
}

// CancelTask changes task state to canceled.
func (s *TaskStore) CancelTask(ctx context.Context, taskID string) error {
	return s.updateWithTaskIDConvert(ctx, taskID, func(task storages.Task) (storages.Task, error) {
		if err := task.State.ValidateStateSwitch(models.CancelTaskEvent); err != nil {
			return storages.Task{}, err
		}

		newTask := task
		newTask.State = models.TaskStateCanceled
		return newTask, nil
	})
}

// BackTaskToWork changes task state to "" from canceled.
func (s *TaskStore) BackTaskToWork(ctx context.Context, taskID string) error {
	return s.updateWithTaskIDConvert(ctx, taskID, func(task storages.Task) (storages.Task, error) {
		if err := task.State.ValidateStateSwitch(models.ToWorkTaskEvent); err != nil {
			return storages.Task{}, err
		}

		newTask := task
		newTask.State = models.TaskStateSimple
		return newTask, nil
	})
}

// PostponeTask postpones task to the next sprint.
// It looks like to cancel and move task to the next sprint template.
func (s *TaskStore) PostponeTask(ctx context.Context, taskID string) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.PostponeTask(ctx, taskID, func(task storages.Task) (storages.Task, storages.Task, error) {
			if err := task.State.ValidateStateSwitch(models.PostponeTaskEvent); err != nil {
				return storages.Task{}, storages.Task{}, err
			}

			if task.Burnt == 0 {
				return task, storages.Task{}, nil
			}

			postponedTask := task
			updatedTask := task

			postponedTask.Burnt = 0
			postponedTask.Points = task.Points - task.Burnt
			updatedTask.State = models.TaskStateCanceled

			return postponedTask, updatedTask, nil
		})
	})
}

func (s *TaskStore) updateWithTaskIDConvert(
	ctx context.Context, taskID string, fn storages.UpdateTaskFn,
) error {
	return s.doWithTaskIDConvert(taskID, func(taskID int64) error {
		return s.dbStore.UpdateTask(ctx, taskID, fn)
	})
}

func (s *TaskStore) doWithTaskIDConvert(taskID string, fn func(taskID int64) error) error {
	id, err := strconv.ParseInt(taskID, 16, 64)
	if err != nil {
		return err
	}

	return fn(id)
}
