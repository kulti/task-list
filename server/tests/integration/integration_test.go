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
	apiURL := os.Getenv("TL_PROXY_URL")
	if apiURL == "" {
		apiURL = "http://127.0.0.1:8080"
	}

	s.Init(apiURL)
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
