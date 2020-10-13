// Package distribution in charge of distributing collected transaction fee and inflated token to all validators and delegators.
// To reduce computation stress, a lazy distribution strategy is brought in. lazy means that the benefit won't be paid directly to contributors automatically.
// The contributors are required to explicitly send transactions to withdraw their benefit, otherwise, their benefit will be kept in the global pool.
//
// [More Details](https://www.irisnet.org/docs/features/distribution.html)
//
// As a quick start:
//
// Withdraw rewards to the withdraw-address(default to the delegator address, you can set to another address via [[setWithdrawAddr]])
//
//  client := test.NewClient()
//  rs, err := dts.Distr().SetWithdrawAddr(dts.Sender().String(), baseTx)
//  require.NoError(dts.T(), err)
//  require.NotEmpty(dts.T(), rs.Hash)
//
package distribution
