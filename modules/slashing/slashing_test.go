package slashing_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/sim"
)

type SlashingTestSuite struct {
	suite.Suite
	sim.TestClient
}

func TestSlashingTestSuite(t *testing.T) {
	suite.Run(t, new(SlashingTestSuite))
}

func (sts *SlashingTestSuite) SetupTest() {
	tc := sim.NewClient()
	sts.TestClient = tc
}

func (sts *SlashingTestSuite) TestQueryParams() {
	params, err := sts.Slashing().QueryParams()
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), params)
}

func (sts *SlashingTestSuite) TestQueryValidatorSigningInfo() {
	validators, err := sts.Staking().QueryValidators(1, 10)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), validators)

	signingInfo, err := sts.Slashing().QueryValidatorSigningInfo(validators[0].ConsensusPubkey)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), signingInfo)
}
