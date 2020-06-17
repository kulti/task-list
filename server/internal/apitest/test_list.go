package apitest

import (
	"math/rand"
	"time"

	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

func (s *APISuite) TestEmptyList() {
	s.newSprint()

	s.checkSprintTaskList()
}

func (s *APISuite) TestSortList() {
	s.newSprint()

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

	seed := time.Now().UnixNano()
	s.T().Log("seed:", seed)
	rand.Seed(seed)
	rand.Shuffle(len(createActions), func(i, j int) {
		createActions[i], createActions[j] = createActions[j], createActions[i]
	})

	for _, a := range createActions {
		a()
	}

	s.todoTask(todoTask.Id)
	s.doneTask(doneTask.Id)
	s.cancelTask(cancelTask.Id)

	s.checkSprintTaskList(todoTask, simpleTask, doneTask, cancelTask)
}
