package sprintstore

import (
	"context"
	"fmt"
	"time"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

type dbStore interface {
	NewSprint(ctx context.Context, opts storages.SprintOpts) error
	CreateTask(ctx context.Context, task models.Task, sprintID string) (string, error)
	ListTasks(ctx context.Context, sprintID string) (models.TaskList, error)
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
	return s.dbStore.CreateTask(ctx, task, sprintID)
}

// ListTasks lists tasks in the sprint.
func (s *SprintStore) ListTasks(ctx context.Context, sprintID string) (models.TaskList, error) {
	return s.dbStore.ListTasks(ctx, sprintID)
}
