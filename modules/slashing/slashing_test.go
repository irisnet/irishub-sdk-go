package slashing_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/suite"
)

type SlashingTestSuite struct {
	suite.Suite
	test.MockClient
}

func TestSlashingTestSuite(t *testing.T) {
	suite.Run(t, new(SlashingTestSuite))
}

func (sts *SlashingTestSuite) SetupTest() {
	tc := test.NewMockClient()
	sts.MockClient = tc
}

func (sts *SlashingTestSuite) TestQueryParams() {
	params, err := sts.Slashing().QueryParams()
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), params)
}

func (sts *SlashingTestSuite) TestQueryValidatorSigningInfo() {
	validators, err := sts.Staking().QueryValidators(1, 10)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), validators)

	signingInfo, err := sts.Slashing().QueryValidatorSigningInfo(validators[0].ConsensusPubkey)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), signingInfo)
	require.NotEmpty(sts.T(), signingInfo.IndexOffset)
}
