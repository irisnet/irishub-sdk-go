package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Distribution interface {
	sdk.Module
	QueryRewards(delegator string) (Rewards, sdk.Error)
	SetWithdrawAddr(withdrawAddr string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}

type Rewards struct {
	Total       sdk.Coins           `json:"total"`
	Delegations []DelegationRewards `json:"delegations"`
	Commission  sdk.Coins           `json:"commission"`
}

type DelegationRewards struct {
	Validator string    `json:"validator"`
	Reward    sdk.Coins `json:"reward"`
}
