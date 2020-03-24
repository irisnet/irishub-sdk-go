package slashing_test

import (
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
	sts.NoError(err)
	sts.NotEmpty(params)
}

func (sts *SlashingTestSuite) TestQueryValidatorSigningInfo() {
	validators, err := sts.Staking().QueryValidators(1, 10)
	sts.NoError(err)
	sts.NotEmpty(validators)

	signingInfo, err := sts.Slashing().QueryValidatorSigningInfo(validators[0].ConsensusPubkey)
	sts.NoError(err)
	sts.NotEmpty(signingInfo)
}
