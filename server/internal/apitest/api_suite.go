package apitest

import (
	"context"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/suite"

	"github.com/kulti/task-list/server/internal/generated/openapicli"
)

type APISuite struct {
	suite.Suite
	cli         *openapicli.APIClient
	ctx         context.Context
	apiURL      string
	sprintTitle string
}

func (s *APISuite) Init(apiURL string) {
	apiCfg := openapicli.NewConfiguration()
	s.cli = openapicli.NewAPIClient(apiCfg)
	s.cli.ChangeBasePath(apiURL + "/api/v1")
	s.ctx = context.Background()
	s.apiURL = apiURL
	s.sprintTitle = faker.Sentence()
}
