package apitest

func (s *APISuite) TestEmptySprintTemplate() {
	tmpl := s.getSprintTemplate()
	s.Require().Empty(tmpl.Tasks)
}
