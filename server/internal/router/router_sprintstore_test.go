package router_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/apitest"
	"github.com/kulti/task-list/server/internal/generated/openapicli"
	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/router"
	"github.com/kulti/task-list/server/internal/services/sprinttmpl"
	"github.com/kulti/task-list/server/internal/storages/memstore"
)

type RouterSprintStoreSuite struct {
	apitest.APISuiteActions
	srv         *httptest.Server
	mockCtrl    *gomock.Controller
	sprintStore *MockSprintStore
}

func (s *RouterSprintStoreSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.sprintStore = NewMockSprintStore(s.mockCtrl)

	store := memstore.NewTaskStore()
	r := router.New(store, s.sprintStore, sprinttmpl.New(store, nil))
	s.srv = httptest.NewServer(r.RootHandler())

	s.Init(s.srv.URL)
}

func (s *RouterSprintStoreSuite) TearDownTest() {
	s.srv.Close()
	s.mockCtrl.Finish()
}

var errTest = errors.New("test error")

func (s *RouterSprintStoreSuite) TestNewSprintError() {
	begin := time.Date(2029, 8, 21, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	s.sprintStore.EXPECT().NewSprint(gomock.Any(), begin, end).Return(errTest)

	_, resp, err := s.Client().CreateTaskList(context.Background(), openapicli.SprintOpts{
		Begin: begin.Format("2006-01-02"),
		End:   end.Format("2006-01-02"),
	})
	s.Require().Error(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (s *RouterSprintStoreSuite) TestCreateTaskError() {
	s.sprintStore.EXPECT().CreateTask(gomock.Any(), gomock.Any(), "current").Return("", errTest)
	s.CreateSprintTaskWithError(http.StatusInternalServerError)
}

func (s *RouterSprintStoreSuite) TestListTasksError() {
	sprintID := faker.Word()
	s.sprintStore.EXPECT().ListTasks(gomock.Any(), sprintID).Return(models.TaskList{}, errTest)
	_, resp, err := s.Client().GetTaskList(context.Background(), sprintID)
	s.Require().Error(err)
	resp.Body.Close()
	s.Require().Equal(http.StatusInternalServerError, resp.StatusCode)
}

func TestRouterSprintStoreErrors(t *testing.T) {
	suite.Run(t, new(RouterSprintStoreSuite))
}
