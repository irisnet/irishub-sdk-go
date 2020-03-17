package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type StakingTx interface {
	Delegate(valAddr string, amount sdk.DecCoin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	Undelegate(valAddr string, amount sdk.DecCoin, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	Redelegate(srcValAddr,
		dstValAddr string,
		amount sdk.DecCoin,
		baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
}

type StakingQueries interface {
	QueryDelegation(delAddr, valAddr string) (Delegation, sdk.Error)
	QueryDelegations(delAddr string) (Delegations, sdk.Error)

	QueryUnbondingDelegation(delAddr, valAddr string) (UnbondingDelegation, sdk.Error)
	QueryUnbondingDelegations(delAddr string) (UnbondingDelegations, sdk.Error)

	QueryRedelegation(delAddr, srcValAddr, dstValAddr string) (Redelegation, sdk.Error)
	QueryRedelegations(delAddr string) (Redelegations, sdk.Error)

	QueryDelegationsTo(valAddr string) (Delegations, sdk.Error)
	QueryUnbondingDelegationsFrom(valAddr string) (UnbondingDelegations, sdk.Error)
	QueryRedelegationsFrom(valAddr string) (Redelegations, sdk.Error)

	QueryValidator(address string) (Validator, sdk.Error)
	QueryValidators(page uint64, size uint16) (Validators, sdk.Error)

	QueryPool() (StakePool, sdk.Error)
	QueryParams() (StakeParams, sdk.Error)
}

type StakingSubscriber interface {
	SubscribeValidatorInfoUpdates(validator string,
		callback func(data EventDataMsgEditValidator)) sdk.Subscription
}

type Staking interface {
	sdk.Module
	StakingTx
	StakingQueries
	StakingSubscriber
}

type Delegation struct {
	DelegatorAddr string `json:"delegator_addr"`
	ValidatorAddr string `json:"validator_addr"`
	Shares        string `json:"shares"`
	Height        int64  `json:"height"`
}
type Delegations []Delegation

type UnbondingDelegations []UnbondingDelegation
type UnbondingDelegation struct {
	TxHash         string   `json:"tx_hash"`
	DelegatorAddr  string   `json:"delegator_addr"`
	ValidatorAddr  string   `json:"validator_addr"`
	CreationHeight int64    `json:"creation_height"`
	MinTime        string   `json:"min_time"`
	InitialBalance sdk.Coin `json:"initial_balance"`
	Balance        sdk.Coin `json:"balance"`
}

type Redelegations []Redelegation
type Redelegation struct {
	DelegatorAddr    string   `json:"delegator_addr"`
	ValidatorSrcAddr string   `json:"validator_src_addr"`
	ValidatorDstAddr string   `json:"validator_dst_addr"`
	CreationHeight   int64    `json:"creation_height"`
	MinTime          string   `json:"min_time"`
	InitialBalance   sdk.Coin `json:"initial_balance"`
	Balance          sdk.Coin `json:"balance"`
	SharesSrc        string   `json:"shares_src"`
	SharesDst        string   `json:"shares_dst"`
}

type Validators []Validator
type Validator struct {
	OperatorAddress string      `json:"operator_address"`
	ConsensusPubkey string      `json:"consensus_pubkey"`
	Jailed          bool        `json:"jailed"`
	Status          string      `json:"status"`
	Tokens          string      `json:"tokens"`
	DelegatorShares string      `json:"delegator_shares"`
	Description     Description `json:"description"`
	BondHeight      int64       `json:"bond_height"`
	UnbondingHeight int64       `json:"unbonding_height"`
	UnbondingTime   string      `json:"unbonding_time"`
	Commission      Commission  `json:"commission"`
}

// DelegatorShareExRate gets the exchange rate of tokens over delegator shares.
// UNITS: tokens/delegator-shares
func (v Validator) DelegatorShareExRate() sdk.Dec {
	delegatorShares, err := sdk.NewDecFromStr(v.DelegatorShares)
	if err != nil {
		return sdk.ZeroDec()
	}

	tokens, err := sdk.NewDecFromStr(v.Tokens)
	if err != nil {
		return sdk.ZeroDec()
	}
	if delegatorShares.IsZero() {
		return sdk.OneDec()
	}
	return tokens.Quo(delegatorShares)
}

type Commission struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
	UpdateTime    string `json:"update_time"`
}

type Description struct {
	Moniker  string `json:"moniker"`
	Identity string `json:"identity"`
	Website  string `json:"website"`
	Details  string `json:"details"`
}

type StakePool struct {
	LooseTokens  string `json:"loose_tokens"`
	BondedTokens string `json:"bonded_tokens"`
}

type StakeParams struct {
	UnbondingTime string `json:"unbonding_time"`
	MaxValidators int    `json:"max_validators"`
}

type EventDataMsgEditValidator struct {
	Height         int64       `json:"height"`
	Hash           string      `json:"hash"`
	Description    Description `json:"description"`
	Address        string      `json:"address"`
	CommissionRate string      `json:"commission_rate"`
}
