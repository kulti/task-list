package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/kulti/task-list/internal/models"
	"github.com/kulti/task-list/internal/storages"
)

const (
	sprintListID  = "sprint"
	todoListID    = "todo"
	backlogListID = "backlog"
)

const (
	taskStateTodo = "todo"
)

type listHandler struct {
	store storages.TaskStore
}

func newListHandler(store storages.TaskStore) listHandler {
	return listHandler{
		store: store,
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
	case "take":
		h.handleTakeTaskToList(w, r, listID)
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to parse new sprint body: %v", err)
		return
	}

	h.store.NewSprint(r.Context(), opts.Title)
}

func (h listHandler) handleCreateTaskInList(w http.ResponseWriter, r *http.Request, listID string) {
	task, err := h.parseTask(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to parse task body: %v", err)
		return
	}

	listIDs := []string{listID}
	if listID == todoListID {
		listIDs = append(listIDs, sprintListID)
		task.State = taskStateTodo
	}
	id, err := h.store.CreateTask(r.Context(), task, listIDs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to create task: %v", err)
		return
	}

	task.ID = id
	data, err := json.Marshal(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to encode body: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (h listHandler) handleTakeTaskToList(w http.ResponseWriter, r *http.Request, listID string) {
	taskID, _ := shiftPath(r.URL.Path)
	if taskID == "" {
		http.NotFound(w, r)
		return
	}

	h.store.TakeTaskToList(r.Context(), taskID, listID)
}

func (h listHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request, listID string) {
	taskID, _ := shiftPath(r.URL.Path)
	if taskID == "" {
		http.NotFound(w, r)
		return
	}

	h.store.DeleteTaskFromList(r.Context(), taskID, listID)
	if listID == sprintListID {
		h.store.DeleteTaskFromList(r.Context(), taskID, todoListID)
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
		return models.Task{}, errors.New("missing required argument 'text'")
	}
	if task.Points == 0 {
		return models.Task{}, errors.New("missing required argument 'points'")
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
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to get week task list: %v", err)
	}

	if len(taskList.Tasks) == 0 {
		taskList.Tasks = []models.Task{}
	} else {
		sort.Slice(taskList.Tasks, func(i, j int) bool {
			switch taskList.Tasks[i].State {
			case "", taskStateTodo:
				return taskList.Tasks[j].State != "" && taskList.Tasks[j].State != taskStateTodo
			case "done":
				return taskList.Tasks[j].State == "canceled"
			}
			return false
		})
	}

	data, err := json.Marshal(&taskList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to encode body: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h listHandler) supportedListID(listID string) bool {
	switch listID {
	case sprintListID, todoListID, backlogListID:
		return true
	default:
		return false
	}
}
