package calservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Service struct {
	srv *calendar.Service
	ids CalendarIDs
}

func New(opts Options) (*Service, error) {
	b, err := ioutil.ReadFile(opts.CredentialPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentional file: %w", err)
	}

	jwt, err := google.JWTConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credential file: %w", err)
	}
	client := oauth2.NewClient(context.Background(), jwt.TokenSource(context.Background()))

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Calendar client: %w", err)
	}

	r, err := os.Open(opts.CalendarIDsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read calendar ids file: %w", err)
	}

	d := json.NewDecoder(r)

	var ids CalendarIDs
	if err := d.Decode(&ids); err != nil {
		return nil, fmt.Errorf("failed to parse calendar ids file: %w", err)
	}

	return &Service{
		srv: srv,
		ids: ids,
	}, nil
}

func (s *Service) GetEvents(ctx context.Context, begin, end time.Time) ([]Event, error) {
	var events []Event // nolint: prealloc
	for _, id := range s.ids.IDs {
		evs, err := s.getCalendarEvents(ctx, begin, end, id.ID)
		if err != nil {
			zap.L().Warn("failed to get calendar events - skip it", zap.String("calendar", id.Name), zap.Error(err))
			continue
		}
		events = append(events, evs...)
	}
	return events, nil
}

func (s *Service) getCalendarEvents(ctx context.Context, begin, end time.Time, id string) ([]Event, error) {
	events, err := s.srv.Events.List(id).
		Context(ctx).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(begin.Format(time.RFC3339)).
		TimeMax(end.AddDate(0, 0, 1).Format(time.RFC3339)).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve events: %w", err)
	}

	evs := make([]Event, len(events.Items))
	for i, item := range events.Items {
		date, err := time.Parse(time.RFC3339, item.Start.DateTime)
		if err != nil {
			return nil, fmt.Errorf("failed to parse event time: %w", err)
		}
		evs[i] = Event{Name: item.Summary, Date: date}
	}

	return evs, nil
}
