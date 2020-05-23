package router_test

import (
	"fmt"
	"testing"

	"github.com/kulti/task-list/internal/generated/openapicli"
	"github.com/stretchr/testify/suite"
)

type SprintTestSuite struct {
	RouterTestSuite
}

func (s *SprintTestSuite) TestEmptyList() {
	s.newSprint()
	s.checkTaskList(openapicli.SPRINT)
}

func (s *SprintTestSuite) TestCreateTask() {
	s.newSprint()

	respTask := s.createTask(openapicli.SPRINT, testTask)

	s.checkTaskList(openapicli.SPRINT, respTask)
}

func (s *SprintTestSuite) TestDeleteTask() {
	s.newSprint()

	respTask := s.createTask(openapicli.SPRINT, testTask)
	s.deleteTaskFromList(respTask.Id, openapicli.SPRINT)

	fmt.Println("~~~~", respTask.Id)
	s.checkTaskList(openapicli.SPRINT)
}

func (s *SprintTestSuite) TestDoneTask() {
	s.newSprint()

	respTask := s.createTask(openapicli.SPRINT, testTask)
	s.doneTask(respTask.Id)

	respTask.State = "done"
	s.checkTaskList(openapicli.SPRINT, respTask)
}

func (s *SprintTestSuite) TestCancelTask() {
	s.newSprint()

	respTask := s.createTask(openapicli.SPRINT, testTask)
	s.cancelTask(respTask.Id)

	respTask.State = "canceled"
	s.checkTaskList(openapicli.SPRINT, respTask)
}

func (s *SprintTestSuite) TestBurnPoints() {
	s.newSprint()

	respTask := s.createTask(openapicli.SPRINT, testTask)
	respTask.Burnt = respTask.Points / 2
	s.updateTask(respTask)
	s.checkTaskList(openapicli.SPRINT, respTask)
}

func TestSprintTasks(t *testing.T) {
	suite.Run(t, new(SprintTestSuite))
}
