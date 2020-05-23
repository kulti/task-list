package router_test

import (
	"testing"

	"github.com/kulti/task-list/internal/generated/openapicli"
	"github.com/stretchr/testify/suite"
)

type TodoTestSuite struct {
	RouterTestSuite
}

func (s *TodoTestSuite) TestCreateTask() {
	s.newSprint()

	respTask := s.createTask(openapicli.TODO, testTask)

	s.checkTaskList(openapicli.TODO, respTask)
	s.checkTaskList(openapicli.SPRINT, respTask)
}

func (s *TodoTestSuite) TestDeleteTaskFromSprintList() {
	s.newSprint()

	respTask := s.createTask(openapicli.TODO, testTask)
	s.deleteTaskFromList(respTask.Id, openapicli.SPRINT)

	s.checkTaskList(openapicli.TODO)
	s.checkTaskList(openapicli.SPRINT)
}

func (s *TodoTestSuite) TestDeleteTaskFromTodoList() {
	s.newSprint()

	respTask := s.createTask(openapicli.TODO, testTask)
	s.deleteTaskFromList(respTask.Id, openapicli.TODO)

	s.checkTaskList(openapicli.TODO)
	s.checkTaskList(openapicli.SPRINT, respTask)
}

func (s *TodoTestSuite) TestTakeTask() {
	s.newSprint()

	respTask := s.createTask(openapicli.SPRINT, testTask)
	s.checkTaskList(openapicli.TODO)

	s.takeTaskToList(respTask.Id, openapicli.TODO)
	respTask.State = "todo"
	s.checkTaskList(openapicli.TODO, respTask)
	s.checkTaskList(openapicli.SPRINT, respTask)
}

func (s *TodoTestSuite) TestNewSprintCleanupTodoList() {
	s.newSprint()
	s.createTask(openapicli.TODO, testTask)
	s.newSprint()

	s.checkTaskList(openapicli.TODO)
}

func TestTodoTasks(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}
