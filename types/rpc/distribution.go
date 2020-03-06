package rpc

import "github.com/irisnet/irishub-sdk-go/types"

type Distribution interface {
	types.Module
	QueryRewards(delegator string) (Rewards, error)
	SetWithdrawAddr(withdrawAddr string, baseTx types.BaseTx) (types.Result, error)
	WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx types.BaseTx) (types.Result, error)
}

type Rewards struct {
	Total       types.Coins         `json:"total"`
	Delegations []DelegationRewards `json:"delegations"`
	Commission  types.Coins         `json:"commission"`
}

type DelegationRewards struct {
	Validator string      `json:"validator"`
	Reward    types.Coins `json:"reward"`
}
