package router_test

import "github.com/kulti/task-list/internal/router/openapi_cli"

var testTask = openapi_cli.Task{
	Text:   "test task",
	Points: 7, //nolint:gomnd
}
