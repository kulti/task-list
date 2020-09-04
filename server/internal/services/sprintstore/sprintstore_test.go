package sprintstore_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/services/sprintstore"
	"github.com/kulti/task-list/server/internal/storages"
)

//go:generate mockgen -package sprintstore_test -destination mock_test.go -source sprintstore.go -mock_names dbStore=MockDBStore

var errTest = errors.New("test error")

type taskMatcher struct {
	models.Task
}

func (m taskMatcher) Matches(x interface{}) bool {
	task, ok := x.(storages.Task)
	if !ok {
		return false
	}
	return task.Text == m.Text
}

func (m taskMatcher) String() string {
	return fmt.Sprintf("is equal to %+v", m.Task)
}

type SprintStoreSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	dbStore  *MockDBStore
	store    *sprintstore.SprintStore
}

func (s *SprintStoreSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.dbStore = NewMockDBStore(s.mockCtrl)
	s.store = sprintstore.New(s.dbStore)
}

func (s *SprintStoreSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *SprintStoreSuite) TestNewSprint() {
	begin := time.Date(2020, 5, 26, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 6, 2, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()

	opts := storages.SprintOpts{
		Title: "26.05 - 02.06",
		Begin: begin,
		End:   end,
	}
	s.dbStore.EXPECT().NewSprint(ctx, opts).Return(nil)
	err := s.store.NewSprint(ctx, begin, end)
	s.Require().NoError(err)

	s.dbStore.EXPECT().NewSprint(ctx, gomock.Any()).Return(errTest)
	err = s.store.NewSprint(ctx, begin, end)
	s.Require().Error(err, errTest.Error())
}

func (s *SprintStoreSuite) TestCreateTask() {
	sprintID := faker.Sentence()
	task := models.Task{
		Text: faker.Sentence(),
	}
	ctx := context.Background()
	dbTaskID := rand.Int63()

	s.dbStore.EXPECT().CreateTask(ctx, taskMatcher{task}, sprintID).Return(dbTaskID, nil)
	newTaskID, err := s.store.CreateTask(ctx, task, sprintID)
	s.Require().NoError(err)
	s.Require().Equal(strconv.FormatInt(dbTaskID, 16), newTaskID)

	s.dbStore.EXPECT().CreateTask(ctx, taskMatcher{task}, sprintID).Return(dbTaskID, errTest)
	_, err = s.store.CreateTask(ctx, task, sprintID)
	s.Require().Error(err, errTest.Error())
}

func (s *SprintStoreSuite) TestListTasks() {
	sprintID := faker.Sentence()
	taskList := storages.TaskList{
		Title: faker.Sentence(),
	}
	for i := 0; i < 3; i++ {
		taskList.Tasks = append(taskList.Tasks, storages.Task{
			Text: faker.Sentence(),
		})
	}
	ctx := context.Background()

	s.dbStore.EXPECT().ListTasks(ctx, sprintID).Return(taskList, nil)
	retTaskList, err := s.store.ListTasks(ctx, sprintID)
	s.Require().NoError(err)
	s.Require().Equal(taskList.Title, retTaskList.Title)
	s.Require().Len(retTaskList.Tasks, len(taskList.Tasks))
	for i := range taskList.Tasks {
		s.Require().Equal(taskList.Tasks[i].Text, retTaskList.Tasks[i].Text)
	}

	s.dbStore.EXPECT().ListTasks(ctx, sprintID).Return(taskList, errTest)
	_, err = s.store.ListTasks(ctx, sprintID)
	s.Require().Error(err, errTest.Error())
}

func TestSprintStore(t *testing.T) {
	suite.Run(t, new(SprintStoreSuite))
}
