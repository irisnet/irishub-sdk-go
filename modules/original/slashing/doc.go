// Package slashing provides the ability to access the interface of the IRISHUB slashing module
//
// In Proof-of-Stake blockchain, validators will get block provisions by staking their token.
// But if they failed to keep online, they will be punished by slashing a small portion of their staked tokens.
// The offline validators will be removed from the validator set and put into jail, which means their voting power is zero.
// During the jail period, these nodes are not even validator candidates. Once the jail period ends, they can send [[unjail]] transactions to free themselves and become validator candidates again.
//
// [More Details](https://www.irisnet.org/docs/features/slashing.html)
//
// As a quick start:
//
//	validators, err := sts.Staking().QueryValidators(1, 10)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), validators)
//
//	signingInfo, err := sts.Slashing().QueryValidatorSigningInfo(validators[0].ConsensusPubkey)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), signingInfo)
//
package slashing
