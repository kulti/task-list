package apitest

import (
	"encoding/json"
	"math"
	"math/rand"

	"github.com/bxcodec/faker/v3"

	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

const (
	taskStateDone     = "done"
	taskStateTodo     = "todo"
	taskStateCanceled = "canceled"
)

func (s *APISuiteActions) taskToRespTask(task openapicli.Task) openapicli.RespTask {
	s.T().Helper()
	data, err := json.Marshal(&task)
	s.Require().NoError(err)

	var respTask openapicli.RespTask
	err = json.Unmarshal(data, &respTask)
	s.Require().NoError(err)

	return respTask
}

func (s *APISuiteActions) errBody(err error) string {
	if apiErr, ok := err.(openapicli.GenericOpenAPIError); ok {
		return string(apiErr.Body())
	}
	return ""
}

func (s *APISuiteActions) testTask() openapicli.Task {
	t := openapicli.Task{
		Text:   faker.Sentence(),
		Points: 1 + rand.Int31n(math.MaxInt16-1),
	}
	s.T().Logf("test task: %+v", t)
	return t
}

func (s *APISuiteActions) respTasksToTemplateTasks(tasks []openapicli.RespTask) []openapicli.TaskTemplate {
	tmplTasks := make([]openapicli.TaskTemplate, len(tasks))
	for i := range tasks {
		tmplTasks[i] = openapicli.TaskTemplate{
			Text:   tasks[i].Text,
			Points: tasks[i].Points,
		}
	}
	return tmplTasks
}
