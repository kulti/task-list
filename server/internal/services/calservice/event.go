package calservice

import "time"

// Event reprents a calendar event.
type Event struct {
	Date      time.Time
	StartDate time.Time
	Name      string
}
