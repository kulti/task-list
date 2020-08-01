package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

const (
	sprintListID  = "sprint"
	backlogListID = "backlog"
)

type listHandler struct {
	store       storages.TaskStore
	tmplService SprintTemplateService
}

func newListHandler(store storages.TaskStore, tmplService SprintTemplateService) listHandler {
	return listHandler{
		store:       store,
		tmplService: tmplService,
	}
}

func (h listHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var listID string
	listID, r.URL.Path = shiftPath(r.URL.Path)

	if !h.supportedListID(listID) {
		http.NotFound(w, r)
		return
	}

	var action string
	action, r.URL.Path = shiftPath(r.URL.Path)
	switch action {
	case "":
		h.handleGetTaskList(w, r, listID)
	case "new":
		if listID == sprintListID {
			h.handleCreateSprint(w, r)
		} else {
			http.NotFound(w, r)
		}
	case "add":
		h.handleCreateTaskInList(w, r, listID)
	case "delete":
		h.handleDeleteTask(w, r, listID)
	default:
		http.NotFound(w, r)
	}
}

func (h listHandler) handleCreateSprint(w http.ResponseWriter, r *http.Request) {
	jsDecoder := json.NewDecoder(r.Body)

	var opts models.SprintOpts
	err := jsDecoder.Decode(&opts)
	if err != nil {
		httpBadRequest(w, "failed to parse new sprint body", err)
		return
	}

	begin, err := time.Parse("2006-01-02", opts.Begin)
	if err != nil {
		httpBadRequest(w, "failed to parse begin time", err)
		return
	}

	end, err := time.Parse("2006-01-02", opts.End)
	if err != nil {
		httpBadRequest(w, "failed to parse end time", err)
		return
	}

	opts.Title = fmt.Sprintf("%02d.%02d - %02d.%02d", begin.Day(), begin.Month(),
		end.Day(), end.Month())

	err = h.store.NewSprint(r.Context(), opts)
	if err != nil {
		httpInternalServerError(w, "failed to create new sprint", err)
		return
	}

	tmpl, err := h.tmplService.Get(r.Context(), begin, end)
	if err != nil {
		zap.L().Warn("failed to get sprint template - skip it", zap.Error(err))
	}

	httpJSON(w, &tmpl)
}

func (h listHandler) handleCreateTaskInList(w http.ResponseWriter, r *http.Request, listID string) {
	task, err := h.parseTask(r.Body)
	if err != nil {
		httpBadRequest(w, "failed to parse task body", err)
		return
	}

	id, err := h.store.CreateTask(r.Context(), task, listID)
	if err != nil {
		httpInternalServerError(w, "failed to create task", err)
		return
	}

	task.ID = id
	httpJSON(w, &task)
}

func (h listHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request, listID string) {
	taskID, _ := shiftPath(r.URL.Path)
	if taskID == "" {
		http.NotFound(w, r)
		return
	}

	err := h.store.DeleteTaskFromList(r.Context(), taskID, listID)
	if err != nil {
		httpInternalServerError(w, "failed to delete task from db", err)
		return
	}
}

func (h listHandler) parseTask(r io.Reader) (models.Task, error) {
	jsDecoder := json.NewDecoder(r)

	var task models.Task
	err := jsDecoder.Decode(&task)
	if err != nil {
		return models.Task{}, err
	}
	if task.Text == "" {
		return models.Task{}, errMissingArgText
	}
	if task.Points == 0 {
		return models.Task{}, errMissingArgPoints
	}

	return task, nil
}

func (h listHandler) handleGetTaskList(w http.ResponseWriter, r *http.Request, listID string) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskList, err := h.store.ListTasks(r.Context(), listID)
	if err != nil {
		httpInternalServerError(w, "failed to get task list from db", err)
		return
	}

	if len(taskList.Tasks) == 0 {
		taskList.Tasks = []models.Task{}
	} else {
		sort.Slice(taskList.Tasks, func(i, j int) bool {
			switch taskList.Tasks[i].State {
			case models.TaskStateTodo:
				return taskList.Tasks[j].State != models.TaskStateTodo
			case "":
				return taskList.Tasks[j].State != "" && taskList.Tasks[j].State != models.TaskStateTodo
			case "done":
				return taskList.Tasks[j].State == "canceled"
			}
			return false
		})
	}

	httpJSON(w, &taskList)
}

func (h listHandler) supportedListID(listID string) bool {
	switch listID {
	case sprintListID, backlogListID:
		return true
	default:
		return false
	}
}
