package slashing_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
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
}
