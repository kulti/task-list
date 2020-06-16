package apitest

func (s *APISuite) TestCreateSprintTask() {
	s.newSprint()

	respTask := s.createSprintTask()

	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList()
}

func (s *APISuite) TestCreateTodoTask() {
	s.newSprint()

	respTask := s.createTodoTask()

	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList(respTask)
}

func (s *APISuite) TestDeleteTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.deleteSprintTask(respTask.Id)

	s.checkSprintTaskList()
	s.checkTodoTaskList()
}

func (s *APISuite) TestDeleteTodoTaskFromSprintList() {
	s.newSprint()

	respTask := s.createTodoTask()
	s.deleteSprintTask(respTask.Id)

	s.checkSprintTaskList()
	s.checkTodoTaskList()
}

func (s *APISuite) TestDeleteTodoTaskFromTodoList() {
	s.newSprint()

	respTask := s.createTodoTask()
	s.deleteTodoTask(respTask.Id)

	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList()
}

func (s *APISuite) TestTakeTask() {
	s.newSprint()

	respTask := s.createSprintTask()

	s.takeTaskToTodoList(respTask.Id)
	respTask.State = "todo"
	s.checkSprintTaskList(respTask)
	s.checkTodoTaskList(respTask)
}

func (s *APISuite) TestTodoTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.todoTask(respTask.Id)

	respTask.State = "todo"
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestDoneTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Burnt = respTask.Points
	respTask.State = "done"
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestCancelTask() {
	s.newSprint()

	respTask := s.createSprintTask()
	s.cancelTask(respTask.Id)

	respTask.State = "canceled"
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

	respTask.State = "done"
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
