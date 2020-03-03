package types

type Distribution interface {
	QueryRewards(delegator string) (Rewards, error)
	SetWithdrawAddr(withdrawAddr string, baseTx BaseTx) (Result, error)
	WithdrawRewards(isValidator bool, onlyFromValidator string, baseTx BaseTx) (Result, error)
}

type Rewards struct {
	Total       Coins               `json:"total"`
	Delegations []DelegationRewards `json:"delegations"`
	Commission  Coins               `json:"commission"`
}

type DelegationRewards struct {
	Validator string `json:"validator"`
	Reward    Coins  `json:"reward"`
}
