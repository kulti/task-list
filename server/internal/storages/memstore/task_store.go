package memstore

import (
	"context"
	"strconv"

	"github.com/kulti/task-list/server/internal/models"
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

func (s *TaskStore) NewSprint(_ context.Context, opts models.SprintOpts) error {
	s.tasks[sprintList] = &taskList{title: opts.Title}
	s.tasks[todoList] = &taskList{title: "Todo"}
	return nil
}

func (s *TaskStore) CreateTask(_ context.Context, task models.Task, listID string) (string, error) {
	task.ID = s.nextID()
	storeTask := &storeTask{
		Task:    task,
		listIDs: []string{listID},
	}
	s.tasks[listID].tasks = append(s.tasks[listID].tasks, storeTask)
	return task.ID, nil
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
			err := s.changeTaskState(s.tasks[sprintList].tasks[i], models.DoneTaskEvent)
			if err != nil {
				return err
			}
			s.tasks[sprintList].tasks[i].Burnt = s.tasks[sprintList].tasks[i].Points
			break
		}
	}
	return nil
}

func (s *TaskStore) UndoneTask(_ context.Context, taskID string) error {
	for i, t := range s.tasks[sprintList].tasks {
		if t.ID == taskID {
			return s.changeTaskState(s.tasks[sprintList].tasks[i], models.UndoneTaskEvent)
		}
	}
	return nil
}

func (s *TaskStore) TodoTask(_ context.Context, taskID string) error {
	for i, t := range s.tasks[sprintList].tasks {
		if t.ID == taskID {
			return s.changeTaskState(s.tasks[sprintList].tasks[i], models.TodoTaskEvent)
		}
	}
	return nil
}

func (s *TaskStore) CancelTask(_ context.Context, taskID string) error {
	for i, t := range s.tasks[sprintList].tasks {
		if t.ID == taskID {
			return s.changeTaskState(s.tasks[sprintList].tasks[i], models.CancelTaskEvent)
		}
	}
	return nil
}

func (s *TaskStore) GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error) {
	return models.SprintTemplate{}, nil
}

func (s *TaskStore) nextID() string {
	s.lastID++
	return strconv.FormatInt(s.lastID, 16)
}

func (s *TaskStore) changeTaskState(task *storeTask, event models.SwitchTaskStateEvent) error {
	state, err := task.State.NextState(event)
	if err != nil {
		return err
	}
	task.State = state
	return nil
}
