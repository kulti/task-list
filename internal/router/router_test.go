package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/kulti/task-list/internal/generated/openapicli"
	"github.com/kulti/task-list/internal/router"
	"github.com/kulti/task-list/internal/storages/memstore"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type RouterTestSuite struct {
	suite.Suite
	srv         *httptest.Server
	cli         *openapicli.APIClient
	ctx         context.Context
	sprintTitle string
}

func (s *RouterTestSuite) SetupTest() {
	store := memstore.NewTaskStore()
	// store, err := pgstore.New("postgres://tl_user:password@127.0.0.1:5432/task_list?sslmode=disable")
	// s.Require().NoError(err)
	r := router.New(store)
	s.ctx = context.Background()
	s.srv = httptest.NewServer(r.RootHandler())
	s.connectToAPI()
}

func (s *RouterTestSuite) TearDownTest() {
	s.srv.Close()
}

func (s *RouterTestSuite) connectToAPI() {
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(s.srv.URL + "/api/v1")
}

func (s *RouterTestSuite) newSprint() {
	s.sprintTitle = "test title"
	opts := openapicli.SprintOpts{
		Title: s.sprintTitle,
	}
	resp, err := s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) checkTaskList(listID openapicli.ListId, tasks ...openapicli.RespTask) {
	taskList, resp, err := s.cli.DefaultApi.GetTaskList(s.ctx, listID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	switch listID {
	case openapicli.SPRINT:
		s.Require().Equal(s.sprintTitle, taskList.Title)
	case openapicli.TODO:
		s.Require().Equal("Todo", taskList.Title)
	}

	if len(tasks) != 0 || len(taskList.Tasks) != 0 {
		s.Require().Equal(tasks, taskList.Tasks)
	}
}

func (s *RouterTestSuite) createTask(listID openapicli.ListId, task openapicli.Task) openapicli.RespTask {
	respTask, resp, err := s.cli.DefaultApi.CreateTask(s.ctx, listID, task)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	switch listID {
	case openapicli.SPRINT:
		s.Require().Empty(respTask.State)
	case openapicli.TODO:
		s.Require().Equal("todo", respTask.State)
	}
	s.Require().NotEmpty(respTask.Id)

	expectedRespTask := s.taskToRespTask(task)
	expectedRespTask.Id = respTask.Id
	expectedRespTask.State = respTask.State
	s.Require().Equal(expectedRespTask, respTask)

	return respTask
}

func (s *RouterTestSuite) takeTaskToList(taskID string, listID openapicli.ListId) {
	resp, err := s.cli.DefaultApi.TakeTask(s.ctx, listID, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) deleteTaskFromList(taskID string, listID openapicli.ListId) {
	resp, err := s.cli.DefaultApi.DeleteTask(s.ctx, listID, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) doneTask(taskID string) {
	resp, err := s.cli.DefaultApi.DoneTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) cancelTask(taskID string) {
	resp, err := s.cli.DefaultApi.CancelTask(s.ctx, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) updateTask(task openapicli.RespTask) {
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

func (s *RouterTestSuite) taskToRespTask(task openapicli.Task) openapicli.RespTask {
	data, err := json.Marshal(&task)
	s.Require().NoError(err)

	var respTask openapicli.RespTask
	err = json.Unmarshal(data, &respTask)
	s.Require().NoError(err)

	return respTask
}

func (s *RouterTestSuite) errBody(err error) string {
	if apiErr, ok := err.(openapicli.GenericOpenAPIError); ok {
		return string(apiErr.Body())
	}
	return ""
}
