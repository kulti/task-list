package router

import "errors"

var (
	errMissingArgText   = errors.New("missing required argument 'text'")
	errMissingArgPoints = errors.New("missing required argument 'points'")
)
