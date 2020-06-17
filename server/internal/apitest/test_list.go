package apitest

import (
	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

func (s *APISuite) TestEmptyList() {
	s.newSprint()

	s.checkSprintTaskList()
}

func (s *APISuite) TestSortList() {
	var (
		todoTask   openapicli.RespTask
		simpleTask openapicli.RespTask
		doneTask   openapicli.RespTask
		cancelTask openapicli.RespTask
	)

	createActions := []func(){
		func() { simpleTask = s.createSprintTask() },
		func() { todoTask = s.createSprintTask(); todoTask.State = taskStateTodo },
		func() { cancelTask = s.createSprintTask(); cancelTask.State = taskStateCanceled },
		func() {
			doneTask = s.createSprintTask()
			doneTask.State = taskStateDone
			doneTask.Burnt = doneTask.Points
		},
	}

	permutate(createActions, func(actions []func()) {
		s.newSprint()

		for _, a := range createActions {
			a()
		}

		s.todoTask(todoTask.Id)
		s.doneTask(doneTask.Id)
		s.cancelTask(cancelTask.Id)

		s.checkSprintTaskList(todoTask, simpleTask, doneTask, cancelTask)
	}, 0)
}

func permutate(a []func(), f func([]func()), i int) {
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
