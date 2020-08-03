package apitest

import (
	"net/http"

	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

func (s *APISuite) TestNotFound() {
	paths := []string{
		"/unknown", "/api/v2", "/api/v1/unknown", "/api/v1/list/unknown",
		"/api/v1/list/sprint/unknown", "/api/v1/list/backlog/new", "/api/v1/list/sprint/delete",
		"/api/v1/task", "/api/v1/task/unknown",
	}

	for _, p := range paths {
		p := p
		s.Run(p, func() {
			resp, err := http.Get(s.apiURL + p) //nolint:noctx
			s.Require().NoError(err)
			resp.Body.Close()
			s.Require().Equal(http.StatusNotFound, resp.StatusCode)
		})
	}
}

func (s *APISuite) TestNewSprintInvalidDates() {
	opts := openapicli.SprintOpts{
		Begin: s.sprintDate.Format("invalid date"),
		End:   s.sprintDate.Format("2006-01-02"),
	}
	_, resp, err := s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)

	opts.Begin, opts.End = opts.End, opts.Begin
	_, resp, err = s.cli.DefaultApi.CreateTaskList(s.ctx, opts)
	s.Require().Error(err)
	defer resp.Body.Close()
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}
