package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kulti/task-list/server/internal/apitest"
	"github.com/kulti/task-list/server/internal/router"
	"github.com/kulti/task-list/server/internal/storages/memstore"
	"github.com/stretchr/testify/suite"
)

type RouterTestSuite struct {
	apitest.APISuite
	srv *httptest.Server
}

func (s *RouterTestSuite) SetupTest() {
	r := router.New(memstore.NewTaskStore(), nil)
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

func TestRouter(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}
