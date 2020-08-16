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
	"github.com/kulti/task-list/server/internal/services/sprintstore"
	"github.com/kulti/task-list/server/internal/storages/memstore"
)

type RouterSprintTmplTestSuite struct {
	apitest.APISuiteActions
	srv         *httptest.Server
	mockCtrl    *gomock.Controller
	tmplService *MockSprintTemplateService
}

func (s *RouterSprintTmplTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.tmplService = NewMockSprintTemplateService(s.mockCtrl)

	store := memstore.NewTaskStore()
	sprintStore := sprintstore.New(store)
	r := router.New(store, sprintStore, s.tmplService)
	s.srv = httptest.NewServer(r.RootHandler())

	s.Init(s.srv.URL)
}

func (s *RouterSprintTmplTestSuite) TearDownTest() {
	s.srv.Close()
	s.mockCtrl.Finish()
}

func (s *RouterSprintTmplTestSuite) TestSomeTasks() {
	begin := time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	expectedTmpl := models.SprintTemplate{
		Tasks: []models.TaskTemplate{
			{Text: faker.Sentence(), Points: 0},
			{Text: faker.Sentence(), Points: 2},
		},
	}

	s.tmplService.EXPECT().Get(gomock.Any(), begin, end).Return(expectedTmpl, nil)

	tmpl := s.createTaskList(begin, end)

	s.Require().Equal(expectedTmpl, s.openapiTmplToModels(tmpl))
}

var errGetTml = errors.New("failed to get sprint template")

func (s *RouterSprintTmplTestSuite) TestGetTemplateError() {
	begin := time.Date(2029, 8, 21, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	s.tmplService.EXPECT().Get(gomock.Any(), begin, end).Return(models.SprintTemplate{}, errGetTml)

	tmpl := s.createTaskList(begin, end)

	s.Require().Equal(models.SprintTemplate{}, s.openapiTmplToModels(tmpl))
}

func (s *RouterSprintTmplTestSuite) createTaskList(begin, end time.Time) openapicli.SprintTemplate {
	opts := openapicli.SprintOpts{
		Begin: begin.Format("2006-01-02"),
		End:   end.Format("2006-01-02"),
	}

	tmpl, resp, err := s.Client().CreateTaskList(context.Background(), opts)
	s.Require().NoError(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	return tmpl
}

func (s *RouterSprintTmplTestSuite) openapiTmplToModels(
	tmpl openapicli.SprintTemplate,
) models.SprintTemplate {
	res := models.SprintTemplate{}
	for _, task := range tmpl.Tasks {
		res.Tasks = append(res.Tasks, models.TaskTemplate{
			ID:     task.Id,
			Text:   task.Text,
			Points: task.Points,
		})
	}
	return res
}

func TestRouterWithSprintTemplateService(t *testing.T) {
	suite.Run(t, new(RouterSprintTmplTestSuite))
}
