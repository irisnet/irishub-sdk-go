package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Distribution interface {
	sdk.Module
	QueryRewards(operator string) (Rewards, sdk.Error)
	QueryWithdrawAddr(delegator string) (string, sdk.Error)
	QueryCommission(validator string) (ValidatorAccumulatedCommission, sdk.Error)
	SetWithdrawAddr(withdrawAddr string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}

type Rewards struct {
	Rewards []DelegationsRewards `json:"rewards"`
	Total   sdk.DecCoins         `json:"total"`
}

type DelegationsRewards struct {
	Validator string       `json:"validator"`
	Reward    sdk.DecCoins `json:"reward"`
}

type ValidatorAccumulatedCommission struct {
	Commission sdk.DecCoins `json:"commission"`
}
