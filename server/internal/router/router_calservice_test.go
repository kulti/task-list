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

func (s *RouterCalServiceTestSuite) TestAllDayEvents() {
	begin := time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC)

	events := []calservice.Event{
		{Name: "test event 1", Date: begin.Add(2 * time.Hour * 24)},
		{Name: "test event 2", Date: begin.Add(5 * time.Hour * 24)},
	}

	s.calService.EXPECT().GetEvents(gomock.Any(), begin, end).Return(events, nil)

	tmpl := s.createTaskList(begin, end)

	s.Require().Len(tmpl.Tasks, 2)
	s.Require().Equal("08.07 - "+events[0].Name, tmpl.Tasks[0].Text)
	s.Require().Equal("11.07 - "+events[1].Name, tmpl.Tasks[1].Text)
}

func (s *RouterCalServiceTestSuite) TestAtTimeEvents() {
	begin := time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 11, 19, 0, 0, 0, 0, time.UTC)

	events := []calservice.Event{
		{Name: "test event 1", StartDate: begin.Add(1 * time.Hour * 24).Add(18*time.Hour + 10*time.Minute)},
		{Name: "test event 2", StartDate: begin.Add(3 * time.Hour * 24).Add(7 * time.Hour)},
	}

	s.calService.EXPECT().GetEvents(gomock.Any(), begin, end).Return(events, nil)

	tmpl := s.createTaskList(begin, end)

	s.Require().Len(tmpl.Tasks, 2)
	s.Require().Equal("14.11 - "+events[0].Name+" (18:10)", tmpl.Tasks[0].Text)
	s.Require().Equal("16.11 - "+events[1].Name+" (07:00)", tmpl.Tasks[1].Text)
}

var errCalService = errors.New("calendar service error")

func (s *RouterCalServiceTestSuite) TestCalendarServiceErrorAffectsNothing() {
	s.calService.EXPECT().GetEvents(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errCalService)

	s.NewSprint()
}

func (s *RouterCalServiceTestSuite) createTaskList(begin, end time.Time) openapicli.SprintTemplate {
	opts := openapicli.SprintOpts{
		Title: "cal service sprint",
		Begin: begin.Format("2006-01-02"),
		End:   end.Format("2006-01-02"),
	}

	tmpl, resp, err := s.Client().CreateTaskList(context.Background(), opts)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	return tmpl
}

func TestRouterWithCalendarService(t *testing.T) {
	suite.Run(t, new(RouterCalServiceTestSuite))
}
