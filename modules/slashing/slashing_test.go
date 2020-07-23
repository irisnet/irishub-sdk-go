package slashing_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/suite"
)

type SlashingTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestSlashingTestSuite(t *testing.T) {
	suite.Run(t, new(SlashingTestSuite))
}

func (sts *SlashingTestSuite) SetupTest() {
	tc := test.GetMock()
	sts.MockClient = tc
}

func (sts *SlashingTestSuite) TestQueryParams() {
	params, err := sts.Slashing().QueryParams()
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), params)
}

func (sts *SlashingTestSuite) TestQueryValidatorSigningInfo() {
	signingInfo, err := sts.Slashing().QueryValidatorSigningInfo("icp1zcjduepqngqqcwa3u7f8d0ds8ynds5nlkt7c7djvllf88gs3lprnvkyw4suqncgprh")
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), signingInfo)
	require.NotEmpty(sts.T(), signingInfo.IndexOffset)
}
