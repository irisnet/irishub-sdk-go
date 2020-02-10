package distr

import "github.com/irisnet/irishub-sdk-go/types"

type Distr interface {
	//QueryWithdrawAddress(delAddr string) (types.AccAddress, error)
	//QueryDelegationDistInfo(delAddr, valAddr string) (types.DelegationDistInfo, error)
	//QueryAllDelegationDistInfo(delAddr string) (types.DelegationDistInfo, error)
	//QueryValidatorDistInfo(valAddr string) (types.DelegationDistInfo, error)
	QueryRewards(address string) (types.Rewards, error)
	//SetWithdrawAddress(delAddr, valAddr string, baseTx types.BaseTx) (types.Result, error)
	//WithdrawDelegationRewardsAll(delAddr string, baseTx types.BaseTx) (types.Result, error)
	//WithdrawDelegationReward(delAddr, valAddr string, baseTx types.BaseTx) (types.Result, error)
	//WithdrawValidatorRewardsAll(valAddr string, baseTx types.BaseTx) (types.Result, error)
}

type distrClient struct {
	types.AbstractClient
}

// params for query 'custom/distr/delegation_dist_info', 'custom/distr/all_delegation_dist_info' and 'withdraw_addr'
type QueryDelegatorParams struct {
	DelegatorAddress types.AccAddress `json:"delegator_address"`
}

// params for query 'custom/distr/rewards'
type QueryRewardsParams struct {
	Address types.AccAddress `json:"address"`
}
