package sprinttmpl

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/services/calservice"
)

// Service represents sprint template service.
type Service struct {
	store      Store
	calService CalService
}

// CalService is an interface to get calendar events.
type CalService interface {
	GetEvents(ctx context.Context, begin, end time.Time) ([]calservice.Event, error)
}

// Store is an interface to get sprint template from persistent storage.
type Store interface {
	GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error)
	PopPostponedTasks(ctx context.Context) ([]models.PostponedTask, error)
	SetSprintTemplate(ctx context.Context, tmpl models.SprintTemplate) error
}

// New creates a new instance of sprint template service.
func New(store Store, calService CalService) *Service {
	return &Service{
		store:      store,
		calService: calService,
	}
}

// Get returns a sprint template for begin-end date period.
// Get is not idempotent because of external dependency (google calendar)
// and postponed tasks which returns only on first call.
func (s *Service) Get(ctx context.Context, begin, end time.Time) (models.SprintTemplate, error) {
	tmpl, err := s.store.GetSprintTemplate(ctx)
	if err != nil {
		return models.SprintTemplate{}, err
	}

	postponedTasks, err := s.store.PopPostponedTasks(ctx)
	if err != nil {
		return models.SprintTemplate{}, err
	}

	for _, task := range postponedTasks {
		tmpl.Tasks = append(tmpl.Tasks, models.TaskTemplate{
			Text:   task.Text,
			Points: task.Points,
		})
	}

	s.extendSprintTemplateWithCalendarEvents(ctx, &tmpl, begin, end)

	return tmpl, nil
}

func (s *Service) extendSprintTemplateWithCalendarEvents(ctx context.Context,
	tmpl *models.SprintTemplate, begin, end time.Time,
) {
	if s.calService == nil {
		return
	}

	events, err := s.calService.GetEvents(ctx, begin, end)
	if err != nil {
		zap.L().Warn("failed to get calendar events - skip it", zap.Error(err))
		return
	}

	calendarTasks := make([]models.TaskTemplate, len(events))
	for i, e := range events {
		var taskName string
		if !e.Date.IsZero() {
			taskName = fmt.Sprintf("%02d.%02d - %s", e.Date.Day(), e.Date.Month(), e.Name)
		} else {
			taskName = fmt.Sprintf("%02d.%02d - %s (%02d:%02d)",
				e.StartDate.Day(), e.StartDate.Month(),
				e.Name,
				e.StartDate.Hour(), e.StartDate.Minute())
		}
		calendarTasks[i] = models.TaskTemplate{Text: taskName}
	}

	sort.Slice(calendarTasks, func(i, j int) bool {
		return calendarTasks[i].Text < calendarTasks[j].Text
	})

	tmpl.Tasks = append(tmpl.Tasks, calendarTasks...)
}

// GetNewSprintTemplate returns task from a new sprint template.
func (s *Service) GetNewSprintTemplate(ctx context.Context) (models.SprintTemplate, error) {
	return s.store.GetSprintTemplate(ctx)
}

// SetNewSprintTemplate updates a new sprint template.
func (s *Service) SetNewSprintTemplate(ctx context.Context, tmpl models.SprintTemplate) error {
	return s.store.SetSprintTemplate(ctx, tmpl)
}
