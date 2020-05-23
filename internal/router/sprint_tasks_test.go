package router_test

import (
	"fmt"
	"testing"

	"github.com/kulti/task-list/internal/router/openapi_cli"
	"github.com/stretchr/testify/suite"
)

type SprintTestSuite struct {
	RouterTestSuite
}

func (s *SprintTestSuite) TestEmptyList() {
	s.newSprint()
	s.checkTaskList(openapi_cli.SPRINT)
}

func (s *SprintTestSuite) TestCreateTask() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.SPRINT, testTask)

	s.checkTaskList(openapi_cli.SPRINT, respTask)
}

func (s *SprintTestSuite) TestDeleteTask() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.SPRINT, testTask)
	s.deleteTaskFromList(respTask.Id, openapi_cli.SPRINT)

	fmt.Println("~~~~", respTask.Id)
	s.checkTaskList(openapi_cli.SPRINT)
}

func (s *SprintTestSuite) TestDoneTask() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.SPRINT, testTask)
	s.doneTask(respTask.Id)

	respTask.State = "done"
	s.checkTaskList(openapi_cli.SPRINT, respTask)
}

func (s *SprintTestSuite) TestCancelTask() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.SPRINT, testTask)
	s.cancelTask(respTask.Id)

	respTask.State = "canceled"
	s.checkTaskList(openapi_cli.SPRINT, respTask)
}

func (s *SprintTestSuite) TestBurnPoints() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.SPRINT, testTask)
	respTask.Burnt = respTask.Points / 2
	s.updateTask(respTask)
	s.checkTaskList(openapi_cli.SPRINT, respTask)
}

func TestSprintTasks(t *testing.T) {
	suite.Run(t, new(SprintTestSuite))
}
