package calservice

// Options contains calendar service options.
// CredentialPath is a path to google service account credential file.
// CalendarIDsPath is a path to file with calendar ids (see CalendarIDs).
type Options struct {
	CredentialPath  string
	CalendarIDsPath string
}

// CalendarIDs represents json schema of file with calendar ids.
type CalendarIDs struct {
	IDs []CalendarID `json:"calendars"`
}

// CalendarID contains calendar info.
type CalendarID struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
