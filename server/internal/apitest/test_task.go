package apitest

import "net/http"

func (s *APISuite) TestCreateSprintTask() {
	s.NewSprint()

	respTask := s.createSprintTask()

	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestDeleteTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.deleteSprintTask(respTask.Id)

	s.checkSprintTaskList()
}

func (s *APISuite) TestTodoTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.todoTask(respTask.Id)

	respTask.State = taskStateTodo
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestDoneTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Burnt = respTask.Points
	respTask.State = taskStateDone
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestCancelTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.cancelTask(respTask.Id)

	respTask.State = taskStateCanceled
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestBackTaskToWork() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.cancelTask(respTask.Id)
	s.backTaskToWork(respTask.Id)

	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestBackTaskToWorkNonCanceledTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.backTaskToWorkWithError(respTask.Id, http.StatusBadRequest)
}

func (s *APISuite) TestBurnPoints() {
	s.NewSprint()

	respTask := s.createSprintTask()
	respTask.Burnt = respTask.Points / 2 //nolint:gomnd
	s.updateTask(respTask)
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestBurnAllPoints() {
	s.NewSprint()

	respTask := s.createSprintTask()
	respTask.Burnt = respTask.Points
	s.updateTask(respTask)

	respTask.State = taskStateDone
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestUndoneTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Burnt = 0
	s.updateTask(respTask)

	respTask.State = ""
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestCancelTaskThatAlreadyDone() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	s.cancelTaskWithError(respTask.Id, http.StatusBadRequest)
}

func (s *APISuite) TestDoneTaskThatAlreadyCanceled() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.cancelTask(respTask.Id)

	s.doneTaskWithError(respTask.Id, http.StatusBadRequest)
}

func (s *APISuite) TestUpdateDoneTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.doneTask(respTask.Id)

	respTask.Points++
	respTask.Burnt = respTask.Points
	s.updateTask(respTask)

	respTask.State = taskStateDone
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestUpdateTodoTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.todoTask(respTask.Id)

	respTask.Points++
	respTask.Burnt = respTask.Points - 1
	s.updateTask(respTask)

	respTask.State = taskStateTodo
	s.checkSprintTaskList(respTask)
}

func (s *APISuite) TestPostponeTask() {
	s.NewSprint()

	respTask1 := s.createSprintTask()
	respTask2 := s.createSprintTask()
	s.postponeTask(respTask1.Id)

	s.checkSprintTaskList(respTask2)

	s.NewSprint(respTask1)
}

func (s *APISuite) TestPostponeCanceledTask() {
	s.NewSprint()

	respTask := s.createSprintTask()
	s.cancelTask(respTask.Id)
	s.postponeTaskWithError(respTask.Id, http.StatusBadRequest)
}
