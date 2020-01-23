package types

// distribution info for a delegation - used to determine entitled rewards
type DelegationDistInfo struct {
	DelegatorAddr           AccAddress `json:"delegator_addr"`
	ValOperatorAddr         ValAddress `json:"val_operator_addr"`
	DelPoolWithdrawalHeight int64      `json:"del_pool_withdrawal_height"` // last time this delegation withdrew rewards
}
