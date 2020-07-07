package apitest

import (
	"net/http"

	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

func (s *APISuite) newSprint() {
	s.T().Helper()
	opts := openapicli.SprintOpts{
		Title: s.sprintTitle,
	}
	tmpl, resp, err := s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Empty(tmpl.Tasks)
}

func (s *APISuite) checkSprintTaskList(tasks ...openapicli.RespTask) {
	s.T().Helper()
	s.checkTaskList(openapicli.SPRINT, s.sprintTitle, tasks...)
}

func (s *APISuite) checkTaskList(listID openapicli.ListId, listTitle string, tasks ...openapicli.RespTask) {
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

func (s *APISuite) createSprintTask() openapicli.RespTask {
	s.T().Helper()
	return s.createTask(openapicli.SPRINT, s.testTask())
}

func (s *APISuite) createTask(listID openapicli.ListId, task openapicli.Task) openapicli.RespTask {
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

func (s *APISuite) deleteSprintTask(taskID string) {
	s.T().Helper()
	s.deleteTaskFromList(taskID, openapicli.SPRINT)
}

func (s *APISuite) deleteTaskFromList(taskID string, listID openapicli.ListId) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DeleteTask(s.ctx, listID, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuite) todoTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.TodoTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuite) doneTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DoneTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuite) doneTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.DoneTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuite) cancelTask(taskID string) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.CancelTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *APISuite) cancelTaskWithError(taskID string, httpStatus int) {
	s.T().Helper()
	resp, err := s.cli.DefaultApi.CancelTask(s.ctx, taskID)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(httpStatus, resp.StatusCode)
}

func (s *APISuite) updateTask(task openapicli.RespTask) {
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
