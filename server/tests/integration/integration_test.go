// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/kulti/task-list/server/internal/apitest"
	"github.com/stretchr/testify/suite"
)

type IntegrationSuite struct {
	apitest.APISuite
}

func (s *IntegrationSuite) SetupTest() {
	apiURL := os.Getenv("TL_SERVER_URL")
	if apiURL == "" {
		apiURL = "http://127.0.0.1:8097"
	}

	s.Init(apiURL)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
