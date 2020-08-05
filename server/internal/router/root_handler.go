package router

import (
	"net/http"

	"github.com/kulti/task-list/server/internal/storages"
)

type rootHandler struct {
	sprintHandler sprintHandler
	taskHandler   taskHandler
}

func newRootHandler(store storages.TaskStore, tmplService SprintTemplateService) rootHandler {
	return rootHandler{
		sprintHandler: newSprintHandler(store, tmplService),
		taskHandler:   newTaskHandler(store),
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
	case "sprint":
		h.sprintHandler.ServeHTTP(w, r)
	case "task":
		h.taskHandler.ServeHTTP(w, r)
	default:
		http.NotFound(w, r)
	}
}
