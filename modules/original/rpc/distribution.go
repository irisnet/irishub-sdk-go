package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

type Distribution interface {
	original.Module
	QueryRewards(delAddrOrValAddr string) (Rewards, original.Error)
	QueryWithdrawAddr(delAddrOrValAddr string) (string, original.Error)
	QueryCommission(validator string) (ValidatorAccumulatedCommission, original.Error)
	SetWithdrawAddr(withdrawAddr string, baseTx original.BaseTx) (original.ResultTx, original.Error)
	WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx original.BaseTx) (original.ResultTx, original.Error)
}

type Rewards struct {
	Rewards []DelegationsRewards `json:"rewards"`
	Total   original.DecCoins    `json:"total"`
}

type DelegationsRewards struct {
	Validator string            `json:"validator"`
	Reward    original.DecCoins `json:"reward"`
}

type ValidatorAccumulatedCommission struct {
	Commission original.DecCoins `json:"commission"`
}
