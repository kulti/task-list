package storages

import "time"

// SprintOpts represents new sprint options.
type SprintOpts struct {
	Title string
	Begin time.Time
	End   time.Time
}
