package apitest

import (
	"strings"

	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

func (s *APISuite) TestEmptyList() {
	s.NewSprint()

	s.checkSprintTaskList()
}

func (s *APISuite) TestSortList() {
	var (
		todoTask   openapicli.RespTask
		simpleTask openapicli.RespTask
		doneTask   openapicli.RespTask
		cancelTask openapicli.RespTask
	)

	createActions := []permutateAction{
		{"simple", func() { simpleTask = s.createSprintTask() }},
		{"todo", func() { todoTask = s.createSprintTask(); todoTask.State = taskStateTodo }},
		{"canceled", func() {
			cancelTask = s.createSprintTask()
			cancelTask.State = taskStateCanceled
		}},
		{"done", func() {
			doneTask = s.createSprintTask()
			doneTask.State = taskStateDone
			doneTask.Burnt = doneTask.Points
		}},
	}

	permutate(createActions, func(actions []permutateAction) {
		s.Run(permutateName(actions), func() {
			s.NewSprint()

			for _, a := range createActions {
				a.fn()
			}

			s.todoTask(todoTask.Id)
			s.doneTask(doneTask.Id)
			s.cancelTask(cancelTask.Id)

			s.checkSprintTaskList(todoTask, simpleTask, doneTask, cancelTask)
		})
	}, 0)
}

type permutateAction struct {
	name string
	fn   func()
}

func permutate(a []permutateAction, f func([]permutateAction), i int) {
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

func permutateName(actions []permutateAction) string {
	b := strings.Builder{}
	b.WriteString("+")
	for _, a := range actions {
		b.WriteString(a.name)
		b.WriteString("+")
	}
	return b.String()
}
