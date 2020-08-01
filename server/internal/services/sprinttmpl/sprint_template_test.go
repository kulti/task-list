package sprinttmpl_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/services/calservice"
	"github.com/kulti/task-list/server/internal/services/sprinttmpl"
)

type SprintTemplateTestSuite struct {
	suite.Suite
	mockCtrl   *gomock.Controller
	calService *MockCalService
	store      *MockStore
	tmpl       *sprinttmpl.Service
}

func (s *SprintTemplateTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.calService = NewMockCalService(s.mockCtrl)
	s.store = NewMockStore(s.mockCtrl)
	s.tmpl = sprinttmpl.New(s.store, s.calService)
}

func (s *SprintTemplateTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *SprintTemplateTestSuite) TestSprintTemplateHasSomeTasks() {
	begin := time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	expectedTmpl := models.SprintTemplate{
		Tasks: []models.TaskTemplate{
			{Text: faker.Sentence(), Points: 0},
			{Text: faker.Sentence(), Points: 2},
		},
	}

	ctx := context.Background()
	s.store.EXPECT().GetSprintTemplate(ctx).Return(expectedTmpl, nil)
	s.calService.EXPECT().GetEvents(ctx, begin, end)

	tmpl, err := s.tmpl.Get(ctx, begin, end)
	s.Require().NoError(err)
	s.Equal(expectedTmpl, tmpl)
}

func (s *SprintTemplateTestSuite) TestSprintTemplateError() {
	begin := time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	ctx := context.Background()
	s.store.EXPECT().GetSprintTemplate(ctx).Return(models.SprintTemplate{}, errGetTemplate)

	_, err := s.tmpl.Get(ctx, begin, end)
	s.Require().Equal(errGetTemplate, err)
}

func (s *SprintTemplateTestSuite) TestAllDayEvents() {
	begin := time.Date(2020, 7, 6, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	events := []calservice.Event{
		{Name: faker.Sentence(), Date: begin.Add(2 * time.Hour * 24)},
		{Name: faker.Sentence(), Date: begin.Add(5 * time.Hour * 24)},
	}

	ctx := context.Background()
	s.store.EXPECT().GetSprintTemplate(ctx)
	s.calService.EXPECT().GetEvents(ctx, begin, end).Return(events, nil)

	tmpl, err := s.tmpl.Get(ctx, begin, end)
	s.Require().NoError(err)

	s.Require().Len(tmpl.Tasks, 2)
	s.Require().Equal("08.07 - "+events[0].Name, tmpl.Tasks[0].Text)
	s.Require().Equal("11.07 - "+events[1].Name, tmpl.Tasks[1].Text)
}

func (s *SprintTemplateTestSuite) TestAtTimeEvents() {
	begin := time.Date(2020, 11, 13, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	events := []calservice.Event{
		{Name: faker.Sentence(), StartDate: begin.Add(1 * time.Hour * 24).Add(18*time.Hour + 10*time.Minute)},
		{Name: faker.Sentence(), StartDate: begin.Add(3 * time.Hour * 24).Add(7 * time.Hour)},
	}

	ctx := context.Background()
	s.store.EXPECT().GetSprintTemplate(ctx)
	s.calService.EXPECT().GetEvents(ctx, begin, end).Return(events, nil)

	tmpl, err := s.tmpl.Get(ctx, begin, end)
	s.Require().NoError(err)

	s.Require().Len(tmpl.Tasks, 2)
	s.Require().Equal("14.11 - "+events[0].Name+" (18:10)", tmpl.Tasks[0].Text)
	s.Require().Equal("16.11 - "+events[1].Name+" (07:00)", tmpl.Tasks[1].Text)
}

func (s *SprintTemplateTestSuite) TestCalendarServiceErrorAffectsNothing() {
	begin := time.Date(2020, 3, 31, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	expectedTmpl := models.SprintTemplate{
		Tasks: []models.TaskTemplate{
			{Text: faker.Sentence(), Points: 0},
			{Text: faker.Sentence(), Points: 2},
		},
	}

	ctx := context.Background()
	s.store.EXPECT().GetSprintTemplate(ctx).Return(expectedTmpl, nil)
	s.calService.EXPECT().GetEvents(ctx, begin, end).Return(nil, errCalService)

	tmpl, err := s.tmpl.Get(ctx, begin, end)
	s.Require().NoError(err)
	s.Equal(expectedTmpl, tmpl)
}

func (s *SprintTemplateTestSuite) TestMissingCalendarServiceAffectsNothing() {
	begin := time.Date(2020, 3, 31, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	expectedTmpl := models.SprintTemplate{
		Tasks: []models.TaskTemplate{
			{Text: faker.Sentence(), Points: 0},
			{Text: faker.Sentence(), Points: 2},
		},
	}

	ctx := context.Background()
	s.store.EXPECT().GetSprintTemplate(ctx).Return(expectedTmpl, nil)

	tmplWithouCalService := sprinttmpl.New(s.store, nil)
	tmpl, err := tmplWithouCalService.Get(ctx, begin, end)
	s.Require().NoError(err)
	s.Equal(expectedTmpl, tmpl)
}

func (s *SprintTemplateTestSuite) TestEventsOrder() {
	begin := time.Date(2020, 3, 1, 0, 0, 0, 0, time.UTC)
	end := begin.Add(7 * 24 * time.Hour)

	sprintTmplTasks := models.SprintTemplate{
		Tasks: []models.TaskTemplate{
			{Text: "a is the first", Points: 1},
			{Text: "z is the second", Points: 26},
		},
	}

	events := []calservice.Event{
		{Name: "test event 3", Date: begin.Add(3 * time.Hour * 24)},
		{Name: "test event 1", Date: begin.Add(2 * time.Hour * 24)},
		{Name: "test event 2", Date: begin.Add(3 * time.Hour * 24)},
	}

	expectedTmpl := models.SprintTemplate{
		Tasks: []models.TaskTemplate{
			sprintTmplTasks.Tasks[0],
			sprintTmplTasks.Tasks[1],
			{Text: "03.03 - " + events[1].Name},
			{Text: "04.03 - " + events[2].Name},
			{Text: "04.03 - " + events[0].Name},
		},
	}

	ctx := context.Background()
	s.store.EXPECT().GetSprintTemplate(ctx).Return(sprintTmplTasks, nil)
	s.calService.EXPECT().GetEvents(ctx, begin, end).Return(events, nil)

	tmpl, err := s.tmpl.Get(ctx, begin, end)
	s.Require().NoError(err)
	s.Require().Equal(expectedTmpl, tmpl)
}

func TestSprintTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(SprintTemplateTestSuite))
}

var (
	errGetTemplate = errors.New("store error")
	errCalService  = errors.New("calendar service error")
)
