// +build integration

package integration_test

import (
	"context"
	"net/http"
	"os"

	"github.com/kulti/task-list/internal/generated/openapicli"
	"github.com/stretchr/testify/suite"
)

type OpenapiCliSuite struct {
	suite.Suite
	apiURL      string
	cli         *openapicli.APIClient
	ctx         context.Context
	sprintTitle string
}

func (s *OpenapiCliSuite) SetupTest() {
	s.apiURL = os.Getenv("TL_SERVER_URL")
	s.ctx = context.Background()
	s.connectToAPI()
}

func (s *OpenapiCliSuite) connectToAPI() {
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(s.apiURL + "/api/v1")
}

func (s *OpenapiCliSuite) newSprint() {
	s.sprintTitle = "test title"
	opts := openapicli.SprintOpts{
		Title: s.sprintTitle,
	}
	resp, err := s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().NoError(err, s.errBody(err))
	defer resp.Body.Close()
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *OpenapiCliSuite) checkTaskList(listID openapicli.ListId, tasks ...openapicli.RespTask) {
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

func (s *OpenapiCliSuite) errBody(err error) string {
	if apiErr, ok := err.(openapicli.GenericOpenAPIError); ok {
		return string(apiErr.Body())
	}
	return ""
}
