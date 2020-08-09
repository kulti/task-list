package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
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
	case "todo":
		h.handleTodoTask(w, r, taskID)
	case "done":
		h.handleDoneTask(w, r, taskID)
	case "cancel":
		h.handleCancelTask(w, r, taskID)
	case "towork":
		h.handleBackTaskToWork(w, r, taskID)
	case "delete":
		h.handleDeleteTask(w, r, taskID)
	case "postpone":
		h.handlePostponeTask(w, r, taskID)
	default:
		http.NotFound(w, r)
	}
}

func (h taskHandler) handleUpdateTask(w http.ResponseWriter, r *http.Request, taskID string) {
	jsDecoder := json.NewDecoder(r.Body)

	var opts models.UpdateOptions
	err := jsDecoder.Decode(&opts)
	if err != nil {
		httpBadRequest(w, "failed to parse body", err)
		return
	}

	err = h.store.UpdateTask(r.Context(), taskID, opts)
	if err != nil {
		httpInternalServerError(w, "failed to update task in db", err)
	}

	if opts.Burnt == opts.Points {
		err := h.store.DoneTask(r.Context(), taskID)
		if err != nil {
			httpInternalServerError(w, "failed to update task in db", err)
		}
	} else {
		err := h.store.UndoneTask(r.Context(), taskID)
		if err != nil {
			httpInternalServerError(w, "failed to update task in db", err)
		}
	}
}

func (h taskHandler) handleTodoTask(w http.ResponseWriter, r *http.Request, taskID string) {
	err := h.store.TodoTask(r.Context(), taskID)
	if err != nil {
		httpInternalServerError(w, "failed to update task in db", err)
	}
}

func (h taskHandler) handleDoneTask(w http.ResponseWriter, r *http.Request, taskID string) {
	err := h.store.DoneTask(r.Context(), taskID)
	if errors.As(err, &models.StateInconsistencyErr{}) {
		httpBadRequest(w, "failed to update task in db", err)
		return
	}
	if err != nil {
		httpInternalServerError(w, "failed to update task in db", err)
	}
}

func (h taskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request, taskID string) {
	err := h.store.DeleteTask(r.Context(), taskID)
	if err != nil {
		httpInternalServerError(w, "failed to update task in db", err)
	}
}

func (h taskHandler) handleCancelTask(w http.ResponseWriter, r *http.Request, taskID string) {
	err := h.store.CancelTask(r.Context(), taskID)
	if errors.As(err, &models.StateInconsistencyErr{}) {
		httpBadRequest(w, "failed to update task in db", err)
		return
	}
	if err != nil {
		httpInternalServerError(w, "failed to update task in db", err)
	}
}

func (h taskHandler) handleBackTaskToWork(w http.ResponseWriter, r *http.Request, taskID string) {
	err := h.store.BackTaskToWork(r.Context(), taskID)
	if errors.As(err, &models.StateInconsistencyErr{}) {
		httpBadRequest(w, "failed to update task in db", err)
		return
	}
	if err != nil {
		httpInternalServerError(w, "failed to update task in db", err)
	}
}

func (h taskHandler) handlePostponeTask(w http.ResponseWriter, r *http.Request, taskID string) {
	err := h.store.PostponeTask(r.Context(), taskID)
	if errors.As(err, &models.StateInconsistencyErr{}) {
		httpBadRequest(w, "failed to postpone task in db", err)
		return
	}
	if err != nil {
		httpInternalServerError(w, "failed to postpone task in db", err)
	}
}
