package sprintstore_test

import (
	"strings"

	"github.com/bxcodec/faker/v3"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

func (s *SprintStoreSuite) TestSortList() {
	tasks := []models.Task{
		{
			State: models.TaskStateTodo,
			Text:  "todo 1",
		},
		{
			State: models.TaskStateTodo,
			Text:  "todo 2",
		},
		{
			State: models.TaskStateSimple,
			Text:  "simple 1",
		},
		{
			State: models.TaskStateSimple,
			Text:  "simple 2",
		},
		{
			State: models.TaskStateCompleted,
			Text:  "done 1",
		},
		{
			State: models.TaskStateCompleted,
			Text:  "done 2",
		},
		{
			State: models.TaskStateCanceled,
			Text:  "canceled 1",
		},
		{
			State: models.TaskStateCanceled,
			Text:  "canceled 2",
		},
	}

	storageTasks := make([]storages.Task, len(tasks))
	for i, t := range tasks {
		storageTasks[i] = storages.Task{
			State: t.State,
			Text:  t.Text,
		}
	}

	permutate(storageTasks, func(storageTasks []storages.Task) {
		s.Run(permutateName(storageTasks), func() {
			sprintID := faker.Word()
			taskList := storages.TaskList{
				Tasks: storageTasks,
			}
			s.dbStore.EXPECT().ListTasks(s.ctx, sprintID).Return(taskList, nil)
			retTaskList, err := s.store.ListTasks(s.ctx, sprintID)
			s.Require().NoError(err)
			s.Require().Len(retTaskList.Tasks, len(tasks))
			for i := range tasks {
				s.Require().Equal(tasks[i].Text, retTaskList.Tasks[i].Text)
			}
		})
	}, 0)
}

func permutate(a []storages.Task, f func([]storages.Task), i int) {
	if i > len(a) {
		f(a)
		return
	}
	permutate(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		permutate(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func permutateName(tasks []storages.Task) string {
	b := strings.Builder{}
	b.WriteString("+")
	for _, t := range tasks {
		b.WriteString(t.Text)
		b.WriteString("+")
	}
	return b.String()
}
