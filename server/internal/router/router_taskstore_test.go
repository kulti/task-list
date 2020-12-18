package router_test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/apitest"
	"github.com/kulti/task-list/server/internal/router"
	"github.com/kulti/task-list/server/internal/services/sprintstore"
	"github.com/kulti/task-list/server/internal/services/sprinttmpl"
	"github.com/kulti/task-list/server/internal/storages/memstore"
)

type RouterTaskStoreSuite struct {
	apitest.APISuiteActions
	srv       *httptest.Server
	mockCtrl  *gomock.Controller
	taskStore *MockTaskStore
	taskID    string
}

func (s *RouterTaskStoreSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.taskStore = NewMockTaskStore(s.mockCtrl)
	s.taskID = strconv.FormatInt(rand.Int63(), 16)

	store := memstore.NewTaskStore()
	sprintStore := sprintstore.New(store)

	r := router.New(s.taskStore, sprintStore, sprinttmpl.New(store, nil))
	s.srv = httptest.NewServer(r.RootHandler())

	s.Init(s.srv.URL)
}

func (s *RouterTaskStoreSuite) TearDownTest() {
	s.srv.Close()
	s.mockCtrl.Finish()
}

func (s *RouterTaskStoreSuite) TestUpdateTaskError() {
	s.taskStore.EXPECT().UpdateTask(gomock.Any(), s.taskID, gomock.Any()).Return(errTest)
	s.UpdateTaskWithError(s.taskID, http.StatusInternalServerError)
}

func (s *RouterTaskStoreSuite) TestTodoTaskError() {
	s.taskStore.EXPECT().TodoTask(gomock.Any(), s.taskID).Return(errTest)
	s.TodoTaskWithError(s.taskID, http.StatusInternalServerError)
}

func (s *RouterTaskStoreSuite) TestDoneTaskError() {
	s.taskStore.EXPECT().DoneTask(gomock.Any(), s.taskID).Return(errTest)
	s.DoneTaskWithError(s.taskID, http.StatusInternalServerError)
}

func (s *RouterTaskStoreSuite) TestDeleteTaskError() {
	s.taskStore.EXPECT().DeleteTask(gomock.Any(), s.taskID).Return(errTest)
	s.DeleteTaskWithError(s.taskID, http.StatusInternalServerError)
}

func (s *RouterTaskStoreSuite) TestCancelTaskError() {
	s.taskStore.EXPECT().CancelTask(gomock.Any(), s.taskID).Return(errTest)
	s.CancelTaskWithError(s.taskID, http.StatusInternalServerError)
}

func (s *RouterTaskStoreSuite) TestBackToWorkTaskError() {
	s.taskStore.EXPECT().BackTaskToWork(gomock.Any(), s.taskID).Return(errTest)
	s.BackTaskToWorkWithError(s.taskID, http.StatusInternalServerError)
}

func (s *RouterTaskStoreSuite) TestPostponeTaskError() {
	s.taskStore.EXPECT().PostponeTask(gomock.Any(), s.taskID).Return(errTest)
	s.PostponeTaskWithError(s.taskID, http.StatusInternalServerError)
}

func TestRouterTaskStoreErrors(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RouterTaskStoreSuite))
}
