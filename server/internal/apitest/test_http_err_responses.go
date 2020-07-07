package apitest

import "net/http"

func (s *APISuite) TestNotFound() {
	paths := []string{"/unknown", "/api/v2", "/api/v1/unknown", "/api/v1/list/unknown",
		"/api/v1/list/sprint/unknown", "/api/v1/list/backlog/new", "/api/v1/list/sprint/delete",
		"/api/v1/task", "/api/v1/task/unknown"}

	//nolint:scopelint
	for _, p := range paths {
		s.Run(p, func() {
			resp, err := http.Get(s.apiURL + p)
			s.Require().NoError(err)
			resp.Body.Close()
			s.Require().Equal(http.StatusNotFound, resp.StatusCode)
		})
	}
}
