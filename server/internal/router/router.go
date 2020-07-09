package router

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/kulti/task-list/server/internal/services/calservice"
	"github.com/kulti/task-list/server/internal/storages"
)

// CalService is an interface to get calendar events.
type CalService interface {
	GetEvents(ctx context.Context, begin, end time.Time) ([]calservice.Event, error)
}

// Router implements TaskListServer interface.
type Router struct {
	rootHandler rootHandler
}

// New returns new instacne of Router.
func New(store storages.TaskStore, calService CalService) *Router {
	return &Router{
		rootHandler: newRootHandler(store, calService),
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
