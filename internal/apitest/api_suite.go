package apitest

import (
	"context"

	"github.com/kulti/task-list/internal/generated/openapicli"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite
	cli         *openapicli.APIClient
	ctx         context.Context
	apiURL      string
	sprintTitle string
}

func (s *APISuite) Init(apiURL string) {
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(apiURL + "/api/v1")
	s.ctx = context.Background()
	s.apiURL = apiURL
	s.sprintTitle = "test title"
}

func (s *APISuite) TestEmptyList() {
	s.newSprint()

	s.checkSprintTaskList()
	s.checkTodoTaskList()
}

func (s *APISuite) TestCreateSprintTask() {
	s.newSprint()

	respTask := s.createSprintTask()

	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList()
}

func (s *APISuite) TestCreateTodoTask() {
	s.newSprint()

	respTask := s.createTodoTask()

	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList(respTask)
}

func (s *APISuite) TestDeleteTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.deleteSprintTask(respTask.Id)

	s.checkSprintTaskList()
	s.checkTodoTaskList()
}

func (s *APISuite) TestDeleteTodoTaskFromSprintList() {
	s.newSprint()

	respTask := s.createTodoTask()
	s.deleteSprintTask(respTask.Id)

	s.checkSprintTaskList()
	s.checkTodoTaskList()
}

func (s *APISuite) TestDeleteTodoTaskFromTodoList() {
	s.newSprint()

	respTask := s.createTodoTask()
	s.deleteTodoTask(respTask.Id)

	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList()
}

func (s *APISuite) TestTakeTask() {
	s.newSprint()

	respTask := s.createSprintTask()

	s.takeTaskToTodoList(respTask.Id)
	respTask.State = "todo"
	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList(respTask)
}

func (s *APISuite) TestNewSprintCleanupTodoList() {
	s.newSprint()

	s.createTodoTask()
	s.newSprint()

	s.checkTodoTaskList()
}

func (s *APISuite) TestDoneTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Burnt = respTask.Points
	respTask.State = "done"
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestCancelTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.cancelTask(respTask.Id)

	respTask.State = "canceled"
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestBurnPoints() {
	s.newSprint()

	respTask := s.createSprintTask()
	respTask.Burnt = respTask.Points / 2 //nolint:gomnd
	s.updateTask(respTask)
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestBurnAllPoints() {
	s.newSprint()

	respTask := s.createSprintTask()
	respTask.Burnt = respTask.Points
	s.updateTask(respTask)

	respTask.State = "done"
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestUndoneTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Burnt = 0
	s.updateTask(respTask)

	respTask.State = ""
	s.checkSprintTaskList(respTask)
}
