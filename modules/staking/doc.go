// Package staking provides staking functionalities for validators and delegators.
// [More Details](https://www.irisnet.org/docs/features/stake.html)
//
// As a quick start:
//
//  baseTx := sdk.BaseTx{
//		From: "test1",
//		Gas:  20000,
//		Memo: "test",
//		Mode: sdk.Commit,
//	}
//
//	//test QueryValidators
//	validators, _ := sts.Staking().QueryValidators(1, 10)
//	validator := validators[0].OperatorAddress
//	amt, _ := sdk.NewIntFromString("20000000000000000000")
//	amount := sdk.NewCoin("iris-atto", amt)
//
//	//test Delegate
//	rs, err := sts.Staking().Delegate(validator, amount, baseTx)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), rs.Hash)
//
//	//test QueryDelegation
//	delegator := sts.Sender().String()
//	d, err := sts.Staking().QueryDelegation(delegator, validator)
//	require.NoError(sts.T(), err)
//	require.Equal(sts.T(), validator, d.ValidatorAddr)
//	require.Equal(sts.T(), delegator, d.DelegatorAddr)
//
//	//test QueryDelegations
//	ds, err := sts.Staking().QueryDelegations(delegator)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), ds)
//
//	//test QueryDelegationsTo
//	ds, err = sts.Staking().QueryDelegationsTo(validator)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), ds)
//
//	//test Undelegate
//	amt, _ = sdk.NewIntFromString("10000000000000000000")
//	amount = sdk.NewCoin("iris-atto", amt)
//	rs, err = sts.Staking().Undelegate(validator, amount, baseTx)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), rs.Hash)
//
//	//test QueryUnbondingDelegation
//	ubd, err := sts.Staking().QueryUnbondingDelegation(delegator, validator)
//	require.NoError(sts.T(), err)
//	require.Equal(sts.T(), validator, ubd.ValidatorAddr)
//	require.Equal(sts.T(), delegator, ubd.DelegatorAddr)
//
//	//test QueryUnbondingDelegations
//	ubds, err := sts.Staking().QueryUnbondingDelegations(delegator)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), ubds)
//
//	//test QueryUnbondingDelegationsFrom
//	uds, err := sts.Staking().QueryUnbondingDelegationsFrom(validator)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), uds)
//
package staking
