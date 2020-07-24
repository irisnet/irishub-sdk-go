package staking_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
)

type StakingTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(StakingTestSuite))
}

func (sts *StakingTestSuite) SetupTest() {
	sts.MockClient = test.GetMock()
}

func (sts *StakingTestSuite) TestQueryValidator() {
	address := "iva13rtezlhpqms02syv27zc0lqc5nt3z4lcnn820h"
	validator, err := sts.Staking().QueryValidator(address)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), address, validator.OperatorAddress)
}

func (sts *StakingTestSuite) TestQueryValidators() {
	validators, err := sts.Staking().QueryValidators(1, 10)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), validators)
}

func (sts *StakingTestSuite) TestQueryDelegationsTo() {
	result, err := sts.Staking().QueryDelegationsTo("iva13rtezlhpqms02syv27zc0lqc5nt3z4lcnn820h")
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), result)
}

func (sts *StakingTestSuite) TestQueryDelegations() {
	result, err := sts.Staking().QueryDelegations("iaa13rtezlhpqms02syv27zc0lqc5nt3z4lcxzd9js")
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), result)
}

func (sts *StakingTestSuite) TestQueryUnbondingDelegationsFrom() {
	result, err := sts.Staking().QueryUnbondingDelegationsFrom("iva1x98k5n7xj0h3udnf5dcdzw85tsfa75qm682jtg")
	fmt.Println(result)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), result)
}

func (sts *StakingTestSuite) TestQueryPool() {
	p, err := sts.Staking().QueryPool()
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), p)
}

func (sts *StakingTestSuite) TestQueryParams() {
	p, err := sts.Staking().QueryParams()
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), p)
}
