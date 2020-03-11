package staking_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
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
		Memo: "test",
		Mode: sdk.Commit,
	}

	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	amt, _ := sdk.NewIntFromString("10000000000000000000")
	amount := sdk.NewCoin("iris-atto", amt)

	rs, err := sts.Staking().Delegate(validator, amount, baseTx)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), rs.Hash)

}

func (sts *StakingTestSuite) TestUndelegate() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Memo: "test",
		Mode: sdk.Commit,
	}

	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	amt, _ := sdk.NewIntFromString("1000000000000000000")
	amount := sdk.NewCoin("iris-atto", amt)

	rs, err := sts.Staking().Undelegate(validator, amount, baseTx)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), rs.Hash)

}

func (sts *StakingTestSuite) TestRedelegate() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Memo: "test",
		Mode: sdk.Commit,
	}

	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	amt, _ := sdk.NewIntFromString("1000000000000000000")
	amount := sdk.NewCoin("iris-atto", amt)

	rs, err := sts.Staking().Undelegate(validator, amount, baseTx)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), rs.Hash)

}

func (sts *StakingTestSuite) TestQueryQueryDelegation() {
	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	delegator := sts.Sender().String()
	d, err := sts.Staking().QueryDelegation(delegator, validator)
	require.True(sts.T(), err.IsNil())
	require.Equal(sts.T(), validator, d.ValidatorAddr)
	require.Equal(sts.T(), delegator, d.DelegatorAddr)
}

func (sts *StakingTestSuite) TestQueryQueryDelegations() {
	delegator := sts.Sender().String()
	ds, err := sts.Staking().QueryDelegations(delegator)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), ds)
}

func (sts *StakingTestSuite) TestQueryUnbondingDelegation() {
	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	delegator := sts.Sender().String()
	ubd, err := sts.Staking().QueryUnbondingDelegation(delegator, validator)
	require.True(sts.T(), err.IsNil())
	require.Equal(sts.T(), validator, ubd.ValidatorAddr)
	require.Equal(sts.T(), delegator, ubd.DelegatorAddr)
}

func (sts *StakingTestSuite) TestQueryUnbondingDelegations() {
	delegator := sts.Sender().String()
	ubds, err := sts.Staking().QueryUnbondingDelegations(delegator)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), ubds)
}

func (sts *StakingTestSuite) TestQueryDelegationsTo() {
	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	ds, err := sts.Staking().QueryDelegationsTo(validator)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), ds)
}

func (sts *StakingTestSuite) TestQueryUnbondingDelegationsFrom() {
	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	ds, err := sts.Staking().QueryUnbondingDelegationsFrom(validator)
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), ds)
}

func (sts *StakingTestSuite) TestQueryValidator() {
	validators, _ := sts.Staking().QueryValidators(1, 10)
	validator := validators[0].OperatorAddress
	v, err := sts.Staking().QueryValidator(validator)
	require.True(sts.T(), err.IsNil())
	require.EqualValues(sts.T(), validators[0], v)
}

func (sts *StakingTestSuite) TestQueryPool() {
	p, err := sts.Staking().QueryPool()
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), p)
}

func (sts *StakingTestSuite) TestQueryParams() {
	p, err := sts.Staking().QueryParams()
	require.True(sts.T(), err.IsNil())
	require.NotEmpty(sts.T(), p)
}
