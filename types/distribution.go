package types

// distribution info for a delegation - used to determine entitled rewards
type DelegationDistInfo struct {
	DelegatorAddr           AccAddress `json:"delegator_addr"`
	ValOperatorAddr         ValAddress `json:"val_operator_addr"`
	DelPoolWithdrawalHeight int64      `json:"del_pool_withdrawal_height"` // last time this delegation withdrew rewards
}

type DelegationsReward struct {
	Validator ValAddress `json:"validator"`
	Reward    Coins      `json:"reward"`
}

type Rewards struct {
	Total       Coins               `json:"total"`
	Delegations []DelegationsReward `json:"delegations"`
	Commission  Coins               `json:"commission"`
}
