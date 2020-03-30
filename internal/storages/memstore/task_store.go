package memstore

import (
	"context"
	"strconv"

	"github.com/kulti/task-list/internal/models"
)

const (
	sprintList = "sprint"
	todoList   = "todo"
)

type TaskStore struct {
	tasks  map[string]*taskList
	lastID int64
}

type taskList struct {
	title string
	tasks storeTasks
}

type storeTasks []*storeTask

func (st storeTasks) toModelTasks() []models.Task {
	mt := make([]models.Task, len(st))
	for i, t := range st {
		mt[i] = t.Task
	}
	return mt
}

type storeTask struct {
	models.Task
	listIDs []string
}

func NewTaskStore() *TaskStore {
	ts := &TaskStore{
		tasks: make(map[string]*taskList),
	}
	ts.tasks[sprintList] = &taskList{}
	ts.tasks[todoList] = &taskList{}
	return ts
}

func (s *TaskStore) NewSprint(_ context.Context, title string) error {
	s.tasks[sprintList] = &taskList{title: title}
	s.tasks[todoList] = &taskList{title: "Todo"}
	return nil
}

func (s *TaskStore) CreateTask(_ context.Context, task models.Task, listIDs []string) (string, error) {
	task.ID = s.nextID()
	storeTask := &storeTask{
		Task:    task,
		listIDs: listIDs,
	}
	for _, listID := range listIDs {
		s.tasks[listID].tasks = append(s.tasks[listID].tasks, storeTask)
	}
	return task.ID, nil
}

func (s *TaskStore) TakeTaskToList(_ context.Context, taskID, listIDs string) error {
	for _, t := range s.tasks[sprintList].tasks {
		if t.ID == taskID {
			t.State = "todo"
			s.tasks[listIDs].tasks = append(s.tasks[listIDs].tasks, t)
			break
		}
	}
	return nil
}

func (s *TaskStore) UpdateTask(ctx context.Context, taskID string, opts models.UpdateOptions) error {
	for _, t := range s.tasks[sprintList].tasks {
		if t.ID == taskID {
			t.Text = opts.Text
			t.Burnt = opts.Burnt
			t.Points = opts.Points
			break
		}
	}
	return nil
}

func (s *TaskStore) DeleteTaskFromList(_ context.Context, taskID, listID string) error {
	for i, t := range s.tasks[listID].tasks {
		if t.ID == taskID {
			s.tasks[listID].tasks = append(s.tasks[listID].tasks[:i], s.tasks[listID].tasks[i+1:]...)
			break
		}
	}
	return nil
}

func (s *TaskStore) ListTasks(_ context.Context, listID string) (models.TaskList, error) {
	l := models.TaskList{
		Title: s.tasks[listID].title,
		Tasks: s.tasks[listID].tasks.toModelTasks(),
	}
	return l, nil
}

func (s *TaskStore) DoneTask(_ context.Context, taskID string) error {
	for i, t := range s.tasks[sprintList].tasks {
		if t.ID == taskID {
			s.tasks[sprintList].tasks[i].State = "done"
			break
		}
	}
	return nil
}

func (s *TaskStore) CancelTask(_ context.Context, taskID string) error {
	for i, t := range s.tasks[sprintList].tasks {
		if t.ID == taskID {
			s.tasks[sprintList].tasks[i].State = "canceled"
			break
		}
	}
	return nil
}

func (s *TaskStore) nextID() string {
	s.lastID++
	return strconv.FormatInt(s.lastID, 16)
}
