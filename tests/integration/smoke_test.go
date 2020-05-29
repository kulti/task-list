// +build integration

package integration_test

import (
	"testing"

	"github.com/kulti/task-list/internal/generated/openapicli"
	"github.com/stretchr/testify/suite"
)

type SmokeSuite struct {
	OpenapiCliSuite
}

func (s *SmokeSuite) TestEmptyList() {
	s.newSprint()
	s.checkTaskList(openapicli.SPRINT)
}

func TestSmoke(t *testing.T) {
	suite.Run(t, new(SmokeSuite))
}
