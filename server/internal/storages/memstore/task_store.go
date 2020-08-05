package memstore

import (
	"context"
	"strconv"

	"github.com/kulti/task-list/server/internal/models"
)

const (
	sprintList = "sprint"
)

type TaskStore struct {
	taskLists      map[string]*taskList
	lastID         int64
	postponedTasks []models.TaskTemplate
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
		taskLists: make(map[string]*taskList),
	}
	ts.taskLists[sprintList] = &taskList{}
	return ts
}

func (s *TaskStore) NewSprint(_ context.Context, opts models.SprintOpts) error {
	s.taskLists[sprintList] = &taskList{title: opts.Title}
	return nil
}

func (s *TaskStore) CreateTask(_ context.Context, task models.Task, _ string) (string, error) {
	task.ID = s.nextID()
	storeTask := &storeTask{
		Task:    task,
		listIDs: []string{sprintList},
	}
	s.taskLists[sprintList].tasks = append(s.taskLists[sprintList].tasks, storeTask)
	return task.ID, nil
}

func (s *TaskStore) UpdateTask(ctx context.Context, taskID string, opts models.UpdateOptions) error {
	for _, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			t.Text = opts.Text
			t.Burnt = opts.Burnt
			t.Points = opts.Points
			break
		}
	}
	return nil
}

func (s *TaskStore) DeleteTask(_ context.Context, taskID string) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			s.taskLists[sprintList].tasks =
				append(s.taskLists[sprintList].tasks[:i], s.taskLists[sprintList].tasks[i+1:]...)
			break
		}
	}
	return nil
}

func (s *TaskStore) ListTasks(_ context.Context, _ string) (models.TaskList, error) {
	l := models.TaskList{
		Title: s.taskLists[sprintList].title,
		Tasks: s.taskLists[sprintList].tasks.toModelTasks(),
	}
	return l, nil
}

func (s *TaskStore) DoneTask(_ context.Context, taskID string) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			err := s.changeTaskState(s.taskLists[sprintList].tasks[i], models.DoneTaskEvent)
			if err != nil {
				return err
			}
			s.taskLists[sprintList].tasks[i].Burnt = s.taskLists[sprintList].tasks[i].Points
			break
		}
	}
	return nil
}

func (s *TaskStore) UndoneTask(_ context.Context, taskID string) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			return s.changeTaskState(s.taskLists[sprintList].tasks[i], models.UndoneTaskEvent)
		}
	}
	return nil
}

func (s *TaskStore) TodoTask(_ context.Context, taskID string) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			return s.changeTaskState(s.taskLists[sprintList].tasks[i], models.TodoTaskEvent)
		}
	}
	return nil
}

func (s *TaskStore) CancelTask(_ context.Context, taskID string) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			return s.changeTaskState(s.taskLists[sprintList].tasks[i], models.CancelTaskEvent)
		}
	}
	return nil
}

func (s *TaskStore) PostponeTask(_ context.Context, taskID string) error {
	for i, t := range s.taskLists[sprintList].tasks {
		if t.ID == taskID {
			_, err := t.State.NextState(models.PostponeTaskEvent)
			if err != nil {
				return err
			}
			s.postponedTasks = append(s.postponedTasks, models.TaskTemplate{Text: t.Text, Points: t.Points})
			s.taskLists[sprintList].tasks = append(s.taskLists[sprintList].tasks[:i],
				s.taskLists[sprintList].tasks[i+1:]...)
			break
		}
	}
	return nil
}

func (s *TaskStore) GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error) {
	tmpl := models.SprintTemplate{
		Tasks: s.postponedTasks,
	}
	s.postponedTasks = nil

	return tmpl, nil
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
