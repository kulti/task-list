package taskstore_test

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/services/taskstore"
	"github.com/kulti/task-list/server/internal/storages"
)

//go:generate mockgen -package taskstore_test -destination mock_test.go -source taskstore.go -mock_names dbStore=MockDBStore

type TaskStoreSuite struct {
	suite.Suite
	mockCtrl  *gomock.Controller
	dbStore   *MockDBStore
	store     *taskstore.TaskStore
	ctx       context.Context
	taskIDStr string
	task      storages.Task
}

func (s *TaskStoreSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.dbStore = NewMockDBStore(s.mockCtrl)
	s.store = taskstore.New(s.dbStore)
	s.ctx = context.Background()

	points := rand.Int31n(math.MaxInt32-10) + 5
	s.task = storages.Task{
		ID:     rand.Int63(),
		Text:   faker.Sentence(),
		State:  models.TaskStateSimple,
		Points: points,
		Burnt:  points / 2,
	}
	s.taskIDStr = strconv.FormatInt(s.task.ID, 16)
}

func (s *TaskStoreSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *TaskStoreSuite) TestDeleteTask() {
	s.dbStore.EXPECT().DeleteTask(s.ctx, s.task.ID)
	s.Require().NoError(s.store.DeleteTask(s.ctx, s.taskIDStr))
}

func (s *TaskStoreSuite) TestTodoTask() {
	expectedTask := s.task
	expectedTask.State = models.TaskStateTodo
	s.checkTaskAfterUpdate(expectedTask, s.store.TodoTask)
}

func (s *TaskStoreSuite) TestTodoTaskInconcistency() {
	s.task.State = models.TaskStateCompleted
	s.checkUpdateInconcistency(s.store.TodoTask)
}

func (s *TaskStoreSuite) TestDoneTask() {
	expectedTask := s.task
	expectedTask.State = models.TaskStateCompleted
	expectedTask.Burnt = expectedTask.Points
	s.checkTaskAfterUpdate(expectedTask, s.store.DoneTask)
}

func (s *TaskStoreSuite) TestDoneTaskInconcistency() {
	s.task.State = models.TaskStateCanceled
	s.checkUpdateInconcistency(s.store.DoneTask)
}

func (s *TaskStoreSuite) TestCancelTask() {
	expectedTask := s.task
	expectedTask.State = models.TaskStateCanceled
	s.checkTaskAfterUpdate(expectedTask, s.store.CancelTask)
}

func (s *TaskStoreSuite) TestCancelTaskInconcistency() {
	s.task.State = models.TaskStateCompleted
	s.checkUpdateInconcistency(s.store.CancelTask)
}

func (s *TaskStoreSuite) TestBackTaskToWork() {
	s.Run("from completed", func() {
		s.task.State = models.TaskStateCompleted
		expectedTask := s.task
		expectedTask.State = models.TaskStateSimple
		s.checkTaskAfterUpdate(expectedTask, s.store.BackTaskToWork)
	})

	s.Run("from canceled", func() {
		s.task.State = models.TaskStateCanceled
		expectedTask := s.task
		expectedTask.State = models.TaskStateSimple
		s.checkTaskAfterUpdate(expectedTask, s.store.BackTaskToWork)
	})
}

func (s *TaskStoreSuite) TestBackTaskToWorkInconcistency() {
	s.task.State = models.TaskStateSimple
	s.checkUpdateInconcistency(s.store.BackTaskToWork)
}

func (s *TaskStoreSuite) TestUpdateTask() {
	s.Run("simple", func() {
		opts := models.UpdateOptions{
			Text:   faker.Sentence(),
			Points: s.task.Points/2 + 10,
			Burnt:  s.task.Burnt / 2,
		}

		expectedTask := s.task
		expectedTask.Text = opts.Text
		expectedTask.Points = opts.Points
		expectedTask.Burnt = opts.Burnt

		s.checkTaskAfterUpdate(expectedTask, func(ctx context.Context, taskID string) error {
			return s.store.UpdateTask(ctx, taskID, opts)
		})
	})

	s.Run("done", func() {
		opts := models.UpdateOptions{
			Text:   faker.Sentence(),
			Points: s.task.Points,
			Burnt:  s.task.Points,
		}

		expectedTask := s.task
		expectedTask.State = models.TaskStateCompleted
		expectedTask.Text = opts.Text
		expectedTask.Points = opts.Points
		expectedTask.Burnt = opts.Burnt

		s.checkTaskAfterUpdate(expectedTask, func(ctx context.Context, taskID string) error {
			return s.store.UpdateTask(ctx, taskID, opts)
		})
	})

	s.Run("undone", func() {
		s.task.State = models.TaskStateCompleted
		s.task.Burnt = s.task.Points

		opts := models.UpdateOptions{
			Text:   faker.Sentence(),
			Points: s.task.Points,
			Burnt:  s.task.Points - 1,
		}

		expectedTask := s.task
		expectedTask.State = models.TaskStateSimple
		expectedTask.Text = opts.Text
		expectedTask.Points = opts.Points
		expectedTask.Burnt = opts.Burnt

		s.checkTaskAfterUpdate(expectedTask, func(ctx context.Context, taskID string) error {
			return s.store.UpdateTask(ctx, taskID, opts)
		})
	})
}

func (s *TaskStoreSuite) TestUpdateTaskInconcistency() {
	s.task.State = models.TaskStateCanceled
	opts := models.UpdateOptions{
		Points: s.task.Points,
		Burnt:  s.task.Points,
	}

	s.checkUpdateInconcistency(func(ctx context.Context, taskID string) error {
		return s.store.UpdateTask(ctx, taskID, opts)
	})
}

func (s *TaskStoreSuite) TestPostpone() {
	s.Run("ok", func() {
		s.task.Burnt = 0
		s.dbStore.EXPECT().
			PostponeTask(s.ctx, s.task.ID, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ int64, fn storages.PostponeTaskFn) error {
				_, ut, err := fn(s.task)
				s.Require().Zero(ut)
				return err
			})

		s.Require().NoError(s.store.PostponeTask(s.ctx, s.taskIDStr))
	})

	s.Run("partially_done", func() {
		s.task.Burnt = s.task.Points - 1
		s.dbStore.EXPECT().
			PostponeTask(s.ctx, s.task.ID, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ int64, fn storages.PostponeTaskFn) error {
				pt, ut, err := fn(s.task)
				s.Require().Equal(s.task.Points-s.task.Burnt, pt.Points)
				s.task.State = models.TaskStateCanceled
				s.Require().Equal(s.task, ut)
				return err
			})

		s.Require().NoError(s.store.PostponeTask(s.ctx, s.taskIDStr))
	})

	s.Run("db_error", func() {
		s.dbStore.EXPECT().
			PostponeTask(s.ctx, s.task.ID, gomock.Any()).Return(errTest)
		s.Require().EqualError(s.store.PostponeTask(s.ctx, s.taskIDStr), errTest.Error())
	})

	s.Run("invalid_task_id", func() {
		s.Require().EqualError(
			s.store.PostponeTask(s.ctx, "invalidID"),
			`strconv.ParseInt: parsing "invalidID": invalid syntax`,
		)
	})

	s.Run("state_inconcictency", func() {
		s.task.State = models.TaskStateCanceled

		s.dbStore.EXPECT().
			PostponeTask(s.ctx, s.task.ID, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ int64, fn storages.PostponeTaskFn) error {
				_, _, err := fn(s.task)
				return err
			})
		err := s.store.PostponeTask(s.ctx, s.taskIDStr)

		s.Require().Error(err)
		s.Require().True(errors.As(err, &models.StateInconsistencyErr{}))
	})
}

func (s *TaskStoreSuite) checkTaskAfterUpdate(
	expectedTask storages.Task, fn func(context.Context, string) error,
) {
	s.Run("ok", func() {
		var updatedTask storages.Task

		s.dbStore.EXPECT().
			UpdateTask(s.ctx, s.task.ID, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ int64, fn storages.UpdateTaskFn) error {
				var err error
				updatedTask, err = fn(s.task)
				return err
			})
		s.Require().NoError(fn(s.ctx, s.taskIDStr))
		s.Require().Equal(expectedTask, updatedTask)
	})

	s.Run("db_error", func() {
		s.dbStore.EXPECT().
			UpdateTask(s.ctx, s.task.ID, gomock.Any()).Return(errTest)
		s.Require().EqualError(fn(s.ctx, s.taskIDStr), errTest.Error())
	})

	s.Run("invalid_task_id", func() {
		s.Require().EqualError(fn(s.ctx, "invalidID"), `strconv.ParseInt: parsing "invalidID": invalid syntax`)
	})
}

func (s *TaskStoreSuite) checkUpdateInconcistency(fn func(context.Context, string) error) {
	s.dbStore.EXPECT().
		UpdateTask(s.ctx, s.task.ID, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ int64, fn storages.UpdateTaskFn) error {
			_, err := fn(s.task)
			return err
		})
	err := fn(s.ctx, s.taskIDStr)

	s.Require().Error(err)
	s.Require().True(errors.As(err, &models.StateInconsistencyErr{}))
}

var errTest = errors.New(faker.Sentence())

func TestTaskStore(t *testing.T) {
	suite.Run(t, new(TaskStoreSuite))
}
