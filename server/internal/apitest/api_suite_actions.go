package apitest

import (
	"context"
	"net/http"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

type APISuiteActions struct {
	suite.Suite
	cli         *openapicli.APIClient
	ctx         context.Context
	apiURL      string
	sprintTitle string
	sprintDate  time.Time
}

func (s *APISuiteActions) Init(apiURL string) {
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(apiURL + "/api/v1")
	s.ctx = context.Background()
	s.apiURL = apiURL
	s.sprintTitle = faker.Sentence()
	s.sprintDate = time.Now()
}

func (s *APISuiteActions) newSprint() {
	s.T().Helper()
	opts := openapicli.SprintOpts{
		Title: s.sprintTitle,
		Begin: s.sprintDate.Format("2006-01-02"),
		End:   s.sprintDate.Format("2006-01-02"),
	}
	tmpl, resp, err := s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Empty(tmpl.Tasks)
}

func (s *APISuiteActions) checkSprintTaskList(tasks ...openapicli.RespTask) {
	s.T().Helper()
	s.checkTaskList(openapicli.SPRINT, s.sprintTitle, tasks...)
}

func (s *APISuiteActions) checkTaskList(
	listID openapicli.ListId, listTitle string, tasks ...openapicli.RespTask,
) {
	s.T().Helper()
	taskList, resp, err := s.cli.DefaultApi.GetTaskList(s.ctx, listID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	s.Require().Equal(listTitle, taskList.Title)

	if len(tasks) != 0 || len(taskList.Tasks) != 0 {
		s.Require().Equal(tasks, taskList.Tasks)
	}
}

func (s *APISuiteActions) createSprintTask() openapicli.RespTask {
	s.T().Helper()
	return s.createTask(openapicli.SPRINT, s.testTask())
}

func (s *APISuiteActions) createTask(listID openapicli.ListId, task openapicli.Task) openapicli.RespTask {
	s.T().Helper()
	respTask, resp, err := s.cli.DefaultApi.CreateTask(s.ctx, listID, task)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	switch listID {
	case openapicli.SPRINT:
		s.Require().Empty(respTask.State)
	default:
		s.Fail("unsupported list id")
	}
	s.Require().NotEmpty(respTask.Id)

	expectedRespTask := s.taskToRespTask(task)
	expectedRespTask.Id = respTask.Id
	expectedRespTask.State = respTask.State
	s.Require().Equal(expectedRespTask, respTask)

	return respTask
}

func (s *APISuiteActions) deleteSprintTask(taskID string) {
	s.T().Helper()
	s.deleteTaskFromList(taskID, openapicli.SPRINT)
}

func (s *APISuiteActions) deleteTaskFromList(taskID string, listID openapicli.ListId) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DeleteTask(s.ctx, listID, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) todoTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.TodoTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) doneTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DoneTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) doneTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DoneTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuiteActions) cancelTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.CancelTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) cancelTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.CancelTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuiteActions) updateTask(task openapicli.RespTask) {
	s.T().Helper()
	opts := openapicli.UpdateOptions{
		Text:   task.Text,
		Burnt:  task.Burnt,
		Points: task.Points,
	}
	resp, err := s.cli.DefaultApi.UpdateTask(s.ctx, task.Id, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}
