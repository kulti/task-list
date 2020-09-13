package memstore

import (
	"context"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

const (
	sprintList = "sprint"
)

type TaskStore struct {
	taskLists      map[string]*taskList
	lastID         int64
	postponedTasks []models.PostponedTask
	sprintTemplate models.SprintTemplate
}

type taskList struct {
	title string
	tasks storeTasks
}

type storeTasks []*storeTask

func (st storeTasks) toModelTasks() []storages.Task {
	mt := make([]storages.Task, len(st))
	for i, t := range st {
		mt[i] = t.Task
	}
	return mt
}

type storeTask struct {
	storages.Task
	listIDs []string
}

func NewTaskStore() *TaskStore {
	ts := &TaskStore{
		taskLists: make(map[string]*taskList),
	}
	ts.taskLists[sprintList] = &taskList{}
	return ts
}

func (s *TaskStore) NewSprint(_ context.Context, opts storages.SprintOpts) error {
	s.taskLists[sprintList] = &taskList{title: opts.Title}
	return nil
}

func (s *TaskStore) CreateTask(_ context.Context, task storages.Task, _ string) (int64, error) {
	task.ID = s.nextID()
	storeTask := &storeTask{
		Task:    task,
		listIDs: []string{sprintList},
	}
	s.taskLists[sprintList].tasks = append(s.taskLists[sprintList].tasks, storeTask)
	return task.ID, nil
}

func (s *TaskStore) UpdateTask(ctx context.Context, taskID int64, fn storages.UpdateTaskFn) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			task, err := fn(t.Task)
			if err != nil {
				return err
			}

			s.taskLists[sprintList].tasks[i].Task = task
			break
		}
	}
	return nil
}

func (s *TaskStore) DeleteTask(_ context.Context, taskID int64) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			s.taskLists[sprintList].tasks =
				append(s.taskLists[sprintList].tasks[:i], s.taskLists[sprintList].tasks[i+1:]...)
			break
		}
	}
	return nil
}

func (s *TaskStore) ListTasks(_ context.Context, _ string) (storages.TaskList, error) {
	l := storages.TaskList{
		Title: s.taskLists[sprintList].title,
		Tasks: s.taskLists[sprintList].tasks.toModelTasks(),
	}
	return l, nil
}

func (s *TaskStore) PostponeTask(_ context.Context, taskID int64, fn storages.PostponeTaskFn) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			postponedTask, updatedTask, err := fn(t.Task)
			if err != nil {
				return err
			}

			s.postponedTasks = append(s.postponedTasks,
				models.PostponedTask{Text: postponedTask.Text, Points: postponedTask.Points})

			if updatedTask.Points != 0 {
				s.taskLists[sprintList].tasks[i].Task = updatedTask
			} else {
				s.taskLists[sprintList].tasks = append(s.taskLists[sprintList].tasks[:i],
					s.taskLists[sprintList].tasks[i+1:]...)
			}
			break
		}
	}
	return nil
}

func (s *TaskStore) PopPostponedTasks(ctx context.Context) ([]models.PostponedTask, error) {
	postponedTasks := s.postponedTasks
	s.postponedTasks = nil
	return postponedTasks, nil
}

func (s *TaskStore) GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error) {
	return s.sprintTemplate, nil
}

func (s *TaskStore) SetSprintTemplate(ctx context.Context, tmpl models.SprintTemplate) error {
	s.sprintTemplate = tmpl
	return nil
}

func (s *TaskStore) nextID() int64 {
	s.lastID++
	return s.lastID
}
