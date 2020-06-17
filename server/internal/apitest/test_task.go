package apitest

func (s *APISuite) TestCreateSprintTask() {
	s.newSprint()

	respTask := s.createSprintTask()

	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestDeleteTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.deleteSprintTask(respTask.Id)

	s.checkSprintTaskList()
}

func (s *APISuite) TestTodoTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.todoTask(respTask.Id)

	respTask.State = taskStateTodo
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestDoneTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Burnt = respTask.Points
	respTask.State = taskStateDone
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestCancelTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.cancelTask(respTask.Id)

	respTask.State = taskStateCanceled
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestBurnPoints() {
	s.newSprint()

	respTask := s.createSprintTask()
	respTask.Burnt = respTask.Points / 2 //nolint:gomnd
	s.updateTask(respTask)
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestBurnAllPoints() {
	s.newSprint()

	respTask := s.createSprintTask()
	respTask.Burnt = respTask.Points
	s.updateTask(respTask)

	respTask.State = taskStateDone
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestUndoneTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Burnt = 0
	s.updateTask(respTask)

	respTask.State = ""
	s.checkSprintTaskList(respTask)
}