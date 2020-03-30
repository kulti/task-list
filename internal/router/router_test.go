package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/kulti/task-list/internal/router"
	"github.com/kulti/task-list/internal/router/openapi_cli"
	"github.com/kulti/task-list/internal/storages/memstore"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/context"
)

type RouterTestSuite struct {
	suite.Suite
	srv         *httptest.Server
	cli         *openapi_cli.APIClient
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
	apiCfg := openapi_cli.NewConfiguration()
	s.cli = openapi_cli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(s.srv.URL + "/api/v1")
}

func (s *RouterTestSuite) newSprint() {
	s.sprintTitle = "test title"
	opts := openapi_cli.SprintOpts{
		Title: s.sprintTitle,
	}
	resp, err := s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) checkTaskList(listID openapi_cli.ListId, tasks ...openapi_cli.RespTask) {
	taskList, resp, err := s.cli.DefaultApi.GetTaskList(s.ctx, listID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	switch listID {
	case openapi_cli.SPRINT:
		s.Require().Equal(s.sprintTitle, taskList.Title)
	case openapi_cli.TODO:
		s.Require().Equal("Todo", taskList.Title)
	}

	if len(tasks) != 0 || len(taskList.Tasks) != 0 {
		s.Require().Equal(tasks, taskList.Tasks)
	}
}

func (s *RouterTestSuite) createTask(listID openapi_cli.ListId, task openapi_cli.Task) openapi_cli.RespTask {
	respTask, resp, err := s.cli.DefaultApi.CreateTask(s.ctx, listID, task)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().Equal("application/json", resp.Header.Get("Content-Type"))

	switch listID {
	case openapi_cli.SPRINT:
		s.Require().Empty(respTask.State)
	case openapi_cli.TODO:
		s.Require().Equal("todo", respTask.State)
	}
	s.Require().NotEmpty(respTask.Id)

	expectedRespTask := s.taskToRespTask(task)
	expectedRespTask.Id = respTask.Id
	expectedRespTask.State = respTask.State
	s.Require().Equal(expectedRespTask, respTask)

	return respTask
}

func (s *RouterTestSuite) takeTaskToList(taskID string, listID openapi_cli.ListId) {
	resp, err := s.cli.DefaultApi.TakeTask(s.ctx, listID, taskID)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) deleteTaskFromList(taskID string, listID openapi_cli.ListId) {
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

func (s *RouterTestSuite) updateTask(task openapi_cli.RespTask) {
	opts := openapi_cli.UpdateOptions{
		Text:   task.Text,
		Burnt:  task.Burnt,
		Points: task.Points,
	}
	resp, err := s.cli.DefaultApi.UpdateTask(s.ctx, task.Id, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *RouterTestSuite) taskToRespTask(task openapi_cli.Task) openapi_cli.RespTask {
	data, err := json.Marshal(&task)
	s.Require().NoError(err)

	var respTask openapi_cli.RespTask
	err = json.Unmarshal(data, &respTask)
	s.Require().NoError(err)

	return respTask
}

func (s *RouterTestSuite) errBody(err error) string {
	if apiErr, ok := err.(openapi_cli.GenericOpenAPIError); ok {
		return string(apiErr.Body())
	}
	return ""
}
