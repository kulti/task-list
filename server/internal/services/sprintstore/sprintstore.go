package sprintstore

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

type dbStore interface {
	NewSprint(ctx context.Context, opts storages.SprintOpts) error
	CreateTask(ctx context.Context, task storages.Task, sprintID string) (int64, error)
	ListTasks(ctx context.Context, sprintID string) (storages.TaskList, error)
}

// SprintStore provides methods to create sprint, create and list tasks in sprint.
type SprintStore struct {
	dbStore dbStore
}

// New creates a new instance of SprintStore.
func New(dbStore dbStore) *SprintStore {
	return &SprintStore{
		dbStore: dbStore,
	}
}

// NewSprint creates a new sprint.
func (s *SprintStore) NewSprint(ctx context.Context, begin, end time.Time) error {
	opts := storages.SprintOpts{
		Begin: begin,
		End:   end,
	}
	opts.Title = fmt.Sprintf("%02d.%02d - %02d.%02d", begin.Day(), begin.Month(),
		end.Day(), end.Month())

	return s.dbStore.NewSprint(ctx, opts)
}

// CreateTask creates a new task in the sprint.
func (s *SprintStore) CreateTask(ctx context.Context, task models.Task, sprintID string) (string, error) {
	newTask := storages.Task{
		Text:   task.Text,
		State:  task.State,
		Points: task.Points,
		Burnt:  task.Burnt,
	}
	newTaskID, err := s.dbStore.CreateTask(ctx, newTask, sprintID)
	return strconv.FormatInt(newTaskID, 16), err
}

// ListTasks lists tasks in the sprint.
func (s *SprintStore) ListTasks(ctx context.Context, sprintID string) (models.TaskList, error) {
	dbTaskList, err := s.dbStore.ListTasks(ctx, sprintID)
	if err != nil {
		return models.TaskList{}, err
	}

	taskList := models.TaskList{
		Title: dbTaskList.Title,
		Tasks: make([]models.Task, len(dbTaskList.Tasks)),
	}
	for i, task := range dbTaskList.Tasks {
		taskList.Tasks[i] = models.Task{
			ID:     strconv.FormatInt(task.ID, 16),
			Text:   task.Text,
			State:  task.State,
			Points: task.Points,
			Burnt:  task.Burnt,
		}
	}

	sort.Slice(taskList.Tasks, func(i, j int) bool {
		otherState := taskList.Tasks[j].State
		if otherState == taskList.Tasks[i].State {
			return taskList.Tasks[i].Text < taskList.Tasks[j].Text
		}

		switch taskList.Tasks[i].State {
		case models.TaskStateTodo:
			return true
		case models.TaskStateSimple:
			return otherState != models.TaskStateTodo
		case models.TaskStateCompleted:
			return otherState == models.TaskStateCanceled
		case models.TaskStateCanceled:
		}
		return false
	})

	return taskList, nil
}
