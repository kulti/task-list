package apitest

func (s *APISuite) TestEmptyList() {
	s.newSprint()

	s.checkSprintTaskList()
}
