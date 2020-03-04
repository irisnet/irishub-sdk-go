package staking_test

import (
	"testing"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StakingTestSuite struct {
	suite.Suite
	sim.TestClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(StakingTestSuite))
}

func (sts *StakingTestSuite) SetupTest() {
	sts.TestClient = sim.NewClient()
}

func (sts *StakingTestSuite) TestDelegate() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	amt, _ := sdk.NewIntFromString("10000000000000000000")
	amount := sdk.NewCoin("iris-atto", amt)

	rs, err := sts.Delegate(validator, amount, baseTx)
	require.NoError(sts.T(), err)
	require.True(sts.T(), rs.IsSuccess())

}

func (sts *StakingTestSuite) TestUndelegate() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	amt, _ := sdk.NewIntFromString("1000000000000000000")
	amount := sdk.NewCoin("iris-atto", amt)

	rs, err := sts.Undelegate(validator, amount, baseTx)
	require.NoError(sts.T(), err)
	require.True(sts.T(), rs.IsSuccess())

}

func (sts *StakingTestSuite) TestRedelegate() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	amt, _ := sdk.NewIntFromString("1000000000000000000")
	amount := sdk.NewCoin("iris-atto", amt)

	rs, err := sts.Undelegate(validator, amount, baseTx)
	require.NoError(sts.T(), err)
	require.True(sts.T(), rs.IsSuccess())

}

func (sts *StakingTestSuite) TestQueryQueryDelegation() {
	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	delegator := sts.Sender().String()
	d, err := sts.QueryDelegation(delegator, validator)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), validator, d.ValidatorAddr)
	require.Equal(sts.T(), delegator, d.DelegatorAddr)
}

func (sts *StakingTestSuite) TestQueryQueryDelegations() {
	delegator := sts.Sender().String()
	ds, err := sts.QueryDelegations(delegator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), ds)
}

func (sts *StakingTestSuite) TestQueryUnbondingDelegation() {
	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	delegator := sts.Sender().String()
	ubd, err := sts.QueryUnbondingDelegation(delegator, validator)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), validator, ubd.ValidatorAddr)
	require.Equal(sts.T(), delegator, ubd.DelegatorAddr)
}

func (sts *StakingTestSuite) TestQueryUnbondingDelegations() {
	delegator := sts.Sender().String()
	ubds, err := sts.QueryUnbondingDelegations(delegator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), ubds)
}

func (sts *StakingTestSuite) TestQueryDelegationsTo() {
	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	ds, err := sts.QueryDelegationsTo(validator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), ds)
}

func (sts *StakingTestSuite) TestQueryUnbondingDelegationsFrom() {
	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	ds, err := sts.QueryUnbondingDelegationsFrom(validator)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), ds)
}

func (sts *StakingTestSuite) TestQueryValidator() {
	validators, _ := sts.QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	v, err := sts.QueryValidator(validator)
	require.NoError(sts.T(), err)
	require.EqualValues(sts.T(), validators[0], v)
}

func (sts *StakingTestSuite) TestQueryPool() {
	p, err := sts.QueryPool()
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), p)
}

func (sts *StakingTestSuite) TestQueryParams() {
	p, err := sts.QueryParams()
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), p)
}
