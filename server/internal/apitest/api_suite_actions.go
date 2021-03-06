package apitest

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	s.sprintDate = time.Now()
}

func (s *APISuiteActions) Client() *openapicli.DefaultApiService {
	return s.cli.DefaultApi
}

func (s *APISuiteActions) NewSprint(tasks ...openapicli.RespTask) {
	s.T().Helper()
	sprintEndData := s.sprintDate.Add(7 * 24 * time.Hour)
	opts := openapicli.SprintOpts{
		Begin: s.sprintDate.Format("2006-01-02"),
		End:   sprintEndData.Format("2006-01-02"),
	}
	s.sprintTitle = fmt.Sprintf("%02d.%02d - %02d.%02d", s.sprintDate.Day(), s.sprintDate.Month(),
		sprintEndData.Day(), sprintEndData.Month())
	tmpl, resp, err := s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	if len(tasks) == 0 {
		s.Require().Empty(tmpl.Tasks)
	} else {
		expectedTasks := s.respTasksToTemplateTasks(tasks)
		s.Require().Equal(expectedTasks, tmpl.Tasks)
	}
}

func (s *APISuiteActions) checkSprintTaskList(tasks ...openapicli.RespTask) {
	s.T().Helper()
	s.checkTaskList(currentSprintID, s.sprintTitle, tasks...)
}

func (s *APISuiteActions) checkTaskList(
	sprintID string, listTitle string, tasks ...openapicli.RespTask,
) {
	s.T().Helper()
	taskList, resp, err := s.cli.DefaultApi.GetTaskList(s.ctx, sprintID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	s.Require().Equal(listTitle, taskList.Title)

	if len(tasks) != 0 || len(taskList.Tasks) != 0 {
		s.Require().Equal(tasks, taskList.Tasks)
	}
}

func (s *APISuiteActions) checkSprintTemplate(tasks ...openapicli.RespTask) {
	s.T().Helper()
	tmpl, resp, err := s.cli.DefaultApi.GetSprintTemplate(s.ctx)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	if len(tasks) != 0 || len(tmpl.Tasks) != 0 {
		s.Require().Equal(s.respTasksToTemplateTasks(tasks), tmpl.Tasks)
	}
}

func (s *APISuiteActions) setSprintTemplate(tasks ...openapicli.RespTask) {
	s.T().Helper()
	tmpl := openapicli.SprintTemplate{
		Tasks: s.respTasksToTemplateTasks(tasks),
	}
	resp, err := s.cli.DefaultApi.SetSprintTemplate(s.ctx, tmpl)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) createSprintTask() openapicli.RespTask {
	s.T().Helper()
	return s.createTask(currentSprintID, s.testTask())
}

func (s *APISuiteActions) CreateSprintTaskWithError(httpStatus int) {
	s.T().Helper()
	_, resp, err := s.cli.DefaultApi.CreateTask(s.ctx, currentSprintID, s.testTask())
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuiteActions) createTask(sprintID string, task openapicli.Task) openapicli.RespTask {
	s.T().Helper()
	respTask, resp, err := s.cli.DefaultApi.CreateTask(s.ctx, sprintID, task)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	s.Require().Empty(respTask.State)
	s.Require().NotEmpty(respTask.Id)

	expectedRespTask := s.taskToRespTask(task)
	expectedRespTask.Id = respTask.Id
	expectedRespTask.State = respTask.State
	s.Require().Equal(expectedRespTask, respTask)

	return respTask
}

func (s *APISuiteActions) deleteTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DeleteTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) DeleteTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DeleteTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuiteActions) todoTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.TodoTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) TodoTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.TodoTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuiteActions) doneTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DoneTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) DoneTaskWithError(taskID string, httpStatus int) {
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

func (s *APISuiteActions) CancelTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.CancelTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuiteActions) backTaskToWork(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.ToworkTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) BackTaskToWorkWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.ToworkTask(s.ctx, taskID)
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

func (s *APISuiteActions) UpdateTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	opts := openapicli.UpdateOptions{}
	resp, err := s.cli.DefaultApi.UpdateTask(s.ctx, taskID, opts)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuiteActions) postponeTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.PostponeTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuiteActions) PostponeTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.PostponeTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}
