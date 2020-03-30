package router

import (
	"net/http"

	"github.com/kulti/task-list/internal/storages"
)

type rootHandler struct {
	listHandler listHandler
	taskHandler taskHandler
}

func newRootHandler(store storages.TaskStore) rootHandler {
	return rootHandler{
		listHandler: newListHandler(store),
		taskHandler: newTaskHandler(store),
	}
}

func (h rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "api":
		h.handleAPI(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h rootHandler) handleAPI(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head != "v1" {
		http.NotFound(w, r)
		return
	}

	head, r.URL.Path = shiftPath(r.URL.Path)
	switch head {
	case "list":
		h.listHandler.ServeHTTP(w, r)
	case "task":
		h.taskHandler.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
