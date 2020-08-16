package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kulti/task-list/server/internal/apitest"
	"github.com/kulti/task-list/server/internal/generated/openapicli"
	"github.com/kulti/task-list/server/internal/router"
	"github.com/kulti/task-list/server/internal/services/sprintstore"
	"github.com/kulti/task-list/server/internal/services/sprinttmpl"
	"github.com/kulti/task-list/server/internal/storages/memstore"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	apitest.APISuite
	srv *httptest.Server
}

func (s *RouterTestSuite) SetupTest() {
	store := memstore.NewTaskStore()
	sprintStore := sprintstore.New(store)
	r := router.New(store, sprintStore, sprinttmpl.New(store, nil))
	s.srv = httptest.NewServer(r.RootHandler())

	s.Init(s.srv.URL)
}

func (s *RouterTestSuite) TearDownTest() {
	s.srv.Close()
}

func (s *RouterTestSuite) TestApiRootNotFound() {
	resp, err := http.Get(s.srv.URL) //nolint:noctx
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *RouterTestSuite) TestNewSprintInvalidJSON() {
	resp, err := http.Post(s.srv.URL+"/api/v1/sprint", "application/json", nil) //nolint:noctx
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *RouterTestSuite) TestCreateTaskInvalidJSON() {
	resp, err := http.Post(s.srv.URL+"/api/v1/sprint/current/add", "application/json", nil) //nolint:noctx
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
	resp, err := http.Post(s.srv.URL+"/api/v1/task/0/update", "application/json", //nolint:noctx
		strings.NewReader("invalid json"))
	s.Require().NoError(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
