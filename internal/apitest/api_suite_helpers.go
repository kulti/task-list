package apitest

import (
	"encoding/json"

	"github.com/kulti/task-list/internal/generated/openapicli"
)

const (
	taskStateDone     = "done"
	taskStateTodo     = "todo"
	taskStateCanceled = "canceled"
)

func (s *APISuite) taskToRespTask(task openapicli.Task) openapicli.RespTask {
	s.T().Helper()
	data, err := json.Marshal(&task)
	s.Require().NoError(err)

	var respTask openapicli.RespTask
	err = json.Unmarshal(data, &respTask)
	s.Require().NoError(err)

	return respTask
}

func (s *APISuite) errBody(err error) string {
	if apiErr, ok := err.(openapicli.GenericOpenAPIError); ok {
		return string(apiErr.Body())
	}
	return ""
}

func (s *APISuite) testTask() openapicli.Task {
	return openapicli.Task{
		Text:   "test task",
		Points: 7, //nolint:gomnd
	}
}
