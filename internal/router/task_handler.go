package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kulti/task-list/internal/models"
	"github.com/kulti/task-list/internal/storages"
)

type taskHandler struct {
	store storages.TaskStore
}

func newTaskHandler(store storages.TaskStore) taskHandler {
	return taskHandler{
		store: store,
	}
}

func (h taskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var taskID string
	taskID, r.URL.Path = shiftPath(r.URL.Path)

	if taskID == "" {
		http.NotFound(w, r)
		return
	}

	var action string
	action, r.URL.Path = shiftPath(r.URL.Path)
	switch action {
	case "update":
		h.handleUpdateTask(w, r, taskID)
	case "done":
		h.handleDoneTask(w, r, taskID)
	case "cancel":
		h.handleCancelTask(w, r, taskID)
	default:
		http.NotFound(w, r)
	}
}

func (h taskHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request, taskID string) {
	jsDecoder := json.NewDecoder(r.Body)

	var opts models.UpdateOptions
	err := jsDecoder.Decode(&opts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to parse new sprint body: %v", err)
		return
	}

	h.store.UpdateTask(r.Context(), taskID, opts)
}

func (h taskHandler) handleDoneTask(_ http.ResponseWriter, r *http.Request, taskID string) {
	h.store.DoneTask(r.Context(), taskID)
}

func (h taskHandler) handleCancelTask(_ http.ResponseWriter, r *http.Request, taskID string) {
	h.store.CancelTask(r.Context(), taskID)
}
