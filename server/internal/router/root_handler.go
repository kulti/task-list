package router

import (
	"net/http"
)

type rootHandler struct {
	sprintHandler sprintHandler
	taskHandler   taskHandler
}

func newRootHandler(taskStore taskStore, sprintStore sprintStore, tmplService sprintTemplateService,
) rootHandler {
	return rootHandler{
		sprintHandler: newSprintHandler(sprintStore, tmplService),
		taskHandler:   newTaskHandler(taskStore),
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
