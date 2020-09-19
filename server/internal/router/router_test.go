package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/apitest"
	"github.com/kulti/task-list/server/internal/generated/openapicli"
	"github.com/kulti/task-list/server/internal/router"
	"github.com/kulti/task-list/server/internal/services/sprintstore"
	"github.com/kulti/task-list/server/internal/services/sprinttmpl"
	"github.com/kulti/task-list/server/internal/services/taskstore"
	"github.com/kulti/task-list/server/internal/storages/memstore"
)

type RouterTestSuite struct {
	apitest.APISuite
	srv *httptest.Server
}

func (s *RouterTestSuite) SetupTest() {
	store := memstore.NewTaskStore()
	taskStore := taskstore.New(store)
	sprintStore := sprintstore.New(store)
	sprinttmplSrv := sprinttmpl.New(store, nil)
	r := router.New(taskStore, sprintStore, sprinttmplSrv)
	s.srv = httptest.NewServer(r.RootHandler())

	s.Init(s.srv.URL)
}

func (s *RouterTestSuite) TearDownTest() {
	s.srv.Close()
}

func (s *RouterTestSuite) TestApiRootNotFound() {
	resp, err := http.Get(s.srv.URL)
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *RouterTestSuite) TestNewSprintInvalidJSON() {
	resp, err := http.Post(s.srv.URL+"/api/v1/sprint", "application/json", nil)
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *RouterTestSuite) TestCreateTaskInvalidJSON() {
	resp, err := http.Post(s.srv.URL+"/api/v1/sprint/current/add", "application/json", nil)
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *RouterTestSuite) TestCreateTaskWithoutText() {
	task := openapicli.Task{
		Points: 10,
	}
	_, resp, err := s.Client().CreateTask(context.Background(), "current", task)
	s.Require().Error(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *RouterTestSuite) TestCreateTaskWithoutPoints() {
	task := openapicli.Task{
		Text: "test text",
	}
	_, resp, err := s.Client().CreateTask(context.Background(), "current", task)
	s.Require().Error(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *RouterTestSuite) TestUpdateTaskInvalidJSON() {
	resp, err := http.Post(s.srv.URL+"/api/v1/task/0/update", "application/json",
		strings.NewReader("invalid json"))
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *RouterTestSuite) TestNewSprintTemplateInvalidJSON() {
	resp, err := http.Post(s.srv.URL+"/api/v1/new_sprint_template", "application/json",
		strings.NewReader("invalid json"))
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *RouterTestSuite) TestNewSprintTemplateMethdoNotAllowed() {
	resp, err := http.Head(s.srv.URL + "/api/v1/new_sprint_template")
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusMethodNotAllowed, resp.StatusCode)
}

func (s *RouterTestSuite) TestCreateSprintMethdoNotAllowed() {
	resp, err := http.Get(s.srv.URL + "/api/v1/sprint")
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusMethodNotAllowed, resp.StatusCode)
}

func (s *RouterTestSuite) TestCreateSprintTaskMethdoNotAllowed() {
	resp, err := http.Get(s.srv.URL + "/api/v1/sprint/anyid/add")
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusMethodNotAllowed, resp.StatusCode)
}

func (s *RouterTestSuite) TestGetSprintListMethdoNotAllowed() {
	resp, err := http.Post(s.srv.URL+"/api/v1/sprint/anyid", "", nil)
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusMethodNotAllowed, resp.StatusCode)
}

func (s *RouterTestSuite) TestTaskActionsMethdoNotAllowed() {
	actions := []string{"update", "todo", "done", "cancel", "towork", "delete", "postpone"}
	for _, action := range actions {
		action := action
		s.Run(action, func() {
			resp, err := http.Get(s.srv.URL + "/api/v1/task/anyid/" + action)
			s.Require().NoError(err)
			resp.Body.Close()
			s.Require().Equal(http.StatusMethodNotAllowed, resp.StatusCode)
		})
	}
}

func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
