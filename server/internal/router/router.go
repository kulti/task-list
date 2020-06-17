package router

import (
	"net/http"

	"github.com/kulti/task-list/server/internal/storages"
	"github.com/rs/cors"
)

// Router implements TaskListServer interface.
type Router struct {
	rootHandler rootHandler
}

// New returns new instacne of Router.
func New(store storages.TaskStore) *Router {
	return &Router{
		rootHandler: newRootHandler(store),
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
