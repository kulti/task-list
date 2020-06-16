package apitest

func (s *APISuite) TestEmptyList() {
	s.newSprint()

	s.checkSprintTaskList()
	s.checkTodoTaskList()
}

func (s *APISuite) TestNewSprintCleanupTodoList() {
	s.newSprint()

	s.createTodoTask()
	s.newSprint()

	s.checkTodoTaskList()
}
