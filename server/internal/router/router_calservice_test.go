package router_test

import (
	context "context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/apitest"
	"github.com/kulti/task-list/server/internal/generated/openapicli"
	"github.com/kulti/task-list/server/internal/router"
	"github.com/kulti/task-list/server/internal/services/calservice"
	"github.com/kulti/task-list/server/internal/storages/memstore"
)

type RouterCalServiceTestSuite struct {
	apitest.APISuiteActions
	srv        *httptest.Server
	mockCtrl   *gomock.Controller
	calService *MockCalService
}

func (s *RouterCalServiceTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.calService = NewMockCalService(s.mockCtrl)

	r := router.New(memstore.NewTaskStore(), s.calService)
	s.srv = httptest.NewServer(r.RootHandler())

	s.Init(s.srv.URL)
}

func (s *RouterCalServiceTestSuite) TearDownTest() {
	s.srv.Close()
	s.mockCtrl.Finish()
}

func (s *RouterCalServiceTestSuite) TestCalendarService() {
	opts := openapicli.SprintOpts{
		Title: "cal service sprint",
		Begin: "2020-07-06",
		End:   "2020-07-12",
	}

	begin, err := time.Parse("2006-01-02", opts.Begin)
	s.Require().NoError(err)

	end, err := time.Parse("2006-01-02", opts.End)
	s.Require().NoError(err)

	events := []calservice.Event{
		{Name: "test event 1", Date: begin.Add(2 * time.Hour * 24)},
		{Name: "test event 2", Date: begin.Add(5 * time.Hour * 24)},
	}

	s.calService.EXPECT().GetEvents(gomock.Any(), begin, end).Return(events, nil)

	tmpl, resp, err := s.Client().CreateTaskList(context.Background(), opts)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	s.Require().Len(tmpl.Tasks, 2)
	s.Require().Equal("08.07 - "+events[0].Name, tmpl.Tasks[0].Text)
	s.Require().Equal("11.07 - "+events[1].Name, tmpl.Tasks[1].Text)
}

var errCalService = errors.New("calendar service error")

func (s *RouterCalServiceTestSuite) TestCalendarServiceErrorAffectsNothing() {
	s.calService.EXPECT().GetEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errCalService)

	s.NewSprint()
}

func TestRouterWithCalendarService(t *testing.T) {
	suite.Run(t, new(RouterCalServiceTestSuite))
}
