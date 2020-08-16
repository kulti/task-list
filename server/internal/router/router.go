package router

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

type sprintTemplateService interface {
	Get(ctx context.Context, begin, end time.Time) (models.SprintTemplate, error)
}

type sprintStore interface {
	NewSprint(ctx context.Context, begin, end time.Time) error
	CreateTask(ctx context.Context, task models.Task, sprintID string) (string, error)
	ListTasks(ctx context.Context, sprintID string) (models.TaskList, error)
}

// Router implements TaskListServer interface.
type Router struct {
	rootHandler rootHandler
}

// New returns new instacne of Router.
func New(store storages.TaskStore, sprintStore sprintStore, tmplService sprintTemplateService) *Router {
	return &Router{
		rootHandler: newRootHandler(store, sprintStore, tmplService),
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
