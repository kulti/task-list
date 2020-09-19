package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kulti/task-list/server/internal/models"
)

type taskHandler struct {
	store taskStore
}

func newTaskHandler(taskStore taskStore) taskHandler {
	return taskHandler{
		store: taskStore,
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
		h.handleWithMethodPost(w, r, taskID, h.handleUpdateTask)
	case "todo":
		h.handleWithMethodPost(w, r, taskID, h.handleTodoTask)
	case "done":
		h.handleWithMethodPost(w, r, taskID, h.handleDoneTask)
	case "cancel":
		h.handleWithMethodPost(w, r, taskID, h.handleCancelTask)
	case "towork":
		h.handleWithMethodPost(w, r, taskID, h.handleBackTaskToWork)
	case "delete":
		h.handleWithMethodPost(w, r, taskID, h.handleDeleteTask)
	case "postpone":
		h.handleWithMethodPost(w, r, taskID, h.handlePostponeTask)
	default:
		http.NotFound(w, r)
	}
}

func (h taskHandler) handleWithMethodPost(
	w http.ResponseWriter, r *http.Request, taskID string,
	fn func(w http.ResponseWriter, r *http.Request, taskID string),
) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fn(w, r, taskID)
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
