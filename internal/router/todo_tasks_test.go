package router_test

import (
	"testing"

	"github.com/kulti/task-list/internal/router/openapi_cli"
	"github.com/stretchr/testify/suite"
)

type TodoTestSuite struct {
	RouterTestSuite
}

func (s *TodoTestSuite) TestCreateTask() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.TODO, testTask)

	s.checkTaskList(openapi_cli.TODO, respTask)
	s.checkTaskList(openapi_cli.SPRINT, respTask)
}

func (s *TodoTestSuite) TestDeleteTaskFromSprintList() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.TODO, testTask)
	s.deleteTaskFromList(respTask.Id, openapi_cli.SPRINT)

	s.checkTaskList(openapi_cli.TODO)
	s.checkTaskList(openapi_cli.SPRINT)
}

func (s *TodoTestSuite) TestDeleteTaskFromTodoList() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.TODO, testTask)
	s.deleteTaskFromList(respTask.Id, openapi_cli.TODO)

	s.checkTaskList(openapi_cli.TODO)
	s.checkTaskList(openapi_cli.SPRINT, respTask)
}

func (s *TodoTestSuite) TestTakeTask() {
	s.newSprint()

	respTask := s.createTask(openapi_cli.SPRINT, testTask)
	s.checkTaskList(openapi_cli.TODO)

	s.takeTaskToList(respTask.Id, openapi_cli.TODO)
	respTask.State = "todo"
	s.checkTaskList(openapi_cli.TODO, respTask)
	s.checkTaskList(openapi_cli.SPRINT, respTask)
}

func (s *TodoTestSuite) TestNewSprintCleanupTodoList() {
	s.newSprint()
	s.createTask(openapi_cli.TODO, testTask)
	s.newSprint()

	s.checkTaskList(openapi_cli.TODO)
}

func TestTodoTasks(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
