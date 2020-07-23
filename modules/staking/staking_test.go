package staking_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/irisnet/irishub-sdk-go/types"

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

func (sts *StakingTestSuite) TestStaking() {
	baseTx := sdk.BaseTx{
		From:     sts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: sts.Account().Password,
	}

	//test QueryValidators
	validators, _ := sts.Staking().QueryValidators(0, 10)
	validator := validators[0].OperatorAddress

	amount, e := sdk.ParseDecCoin("20iris")
	require.NoError(sts.T(), e)

	//test Delegate
	rs, err := sts.Staking().Delegate(validator, amount, baseTx)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), rs.Hash)

	//test QueryDelegation
	delegator := sts.Account().Address.String()
	d, err := sts.Staking().QueryDelegation(delegator, validator)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), validator, d.ValidatorAddr)
	require.Equal(sts.T(), delegator, d.DelegatorAddr)

	//test QueryDelegations
	ds, err := sts.Staking().QueryDelegations(delegator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), ds)

	//test QueryDelegationsTo
	ds, err = sts.Staking().QueryDelegationsTo(validator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), ds)

	//test Undelegate
	amount, e = sdk.ParseDecCoin("10iris")
	require.NoError(sts.T(), e)

	rs, err = sts.Staking().Undelegate(validator, amount, baseTx)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), rs.Hash)

	//test QueryUnbondingDelegation
	ubd, err := sts.Staking().QueryUnbondingDelegation(delegator, validator)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), validator, ubd.ValidatorAddr)
	require.Equal(sts.T(), delegator, ubd.DelegatorAddr)

	//test QueryUnbondingDelegations
	ubds, err := sts.Staking().QueryUnbondingDelegations(delegator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), ubds)

	//test QueryUnbondingDelegationsFrom
	uds, err := sts.Staking().QueryUnbondingDelegationsFrom(validator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), uds)
}

func (sts *StakingTestSuite) TestQueryValidator() {
	validators, err := sts.Staking().QueryValidator("iva13rtezlhpqms02syv27zc0lqc5nt3z4lcnn820h")
	require.NoError(sts.T(), err)
	fmt.Println(validators)
}

func (sts *StakingTestSuite) TestQueryValidators() {
	validators, err := sts.Staking().QueryValidators(1, 10)
	fmt.Println(validators)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), validators)
}

func (sts *StakingTestSuite) TestQueryQueryDelegations() {
	result, err := sts.Staking().QueryDelegationsTo("iva13rtezlhpqms02syv27zc0lqc5nt3z4lcnn820h")
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
