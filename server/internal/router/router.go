package router

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/kulti/task-list/server/internal/models"
)

type sprintTemplateService interface {
	Get(ctx context.Context, begin, end time.Time) (models.SprintTemplate, error)
}

type sprintStore interface {
	NewSprint(ctx context.Context, begin, end time.Time) error
	CreateTask(ctx context.Context, task models.Task, sprintID string) (string, error)
	ListTasks(ctx context.Context, sprintID string) (models.TaskList, error)
}

type taskStore interface {
	DeleteTask(ctx context.Context, taskID string) error
	UpdateTask(ctx context.Context, taskID string, points models.UpdateOptions) error
	TodoTask(ctx context.Context, taskID string) error
	DoneTask(ctx context.Context, taskID string) error
	CancelTask(ctx context.Context, taskID string) error
	BackTaskToWork(ctx context.Context, taskID string) error
	PostponeTask(ctx context.Context, taskID string) error
}

// Router implements TaskListServer interface.
type Router struct {
	rootHandler rootHandler
}

// New returns new instacne of Router.
func New(taskStore taskStore, sprintStore sprintStore, tmplService sprintTemplateService) *Router {
	return &Router{
		rootHandler: newRootHandler(taskStore, sprintStore, tmplService),
	}
}

// RootHandler returns root handler.
func (r *Router) RootHandler() http.Handler {
	c := cors.New(cors.Options{
		// AllowedOrigins: []string{"http://foo.com", "http://foo.com:8080"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})
	return c.Handler(r.rootHandler)
}
