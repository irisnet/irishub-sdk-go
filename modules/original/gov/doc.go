// Package gov in charge of distributing collected transaction fee and inflated token to all validators and delegators.
// To reduce computation stress, a lazy distribution strategy is brought in. lazy means that the benefit won't be paid directly to contributors automatically.
// The contributors are required to explicitly send transactions to withdraw their benefit, otherwise, their benefit will be kept in the global pool.
//
// [More Details](https://www.irisnet.org/docs/features/distribution.html)
//
// As a quick start:
//
// Deposit tokens for an active proposal.
//
//  client := test.NewClient()
//	amt, _ := sdk.NewIntFromString("10000000000000000000000")
//	amount := sdk.NewCoins(sdk.NewCoin("iris-atto", amt))
//	proposalID := uint64(12)
//	rs, err := gts.Gov().Deposit(proposalID, amount, baseTx)
//	require.NoError(gts.T(), err)
//	require.NotEmpty(gts.T(), rs.Hash)
//
package gov
