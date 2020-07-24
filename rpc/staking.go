package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"time"
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
	QueryDelegations(delAddr string) (DelegationResponses, sdk.Error)

	QueryUnbondingDelegation(delAddr, valAddr string) (UnbondingDelegation, sdk.Error)
	QueryUnbondingDelegations(delAddr string) (UnbondingDelegations, sdk.Error)

	QueryRedelegation(delAddr, srcValAddr, dstValAddr string) (Redelegation, sdk.Error)
	QueryRedelegations(delAddr string) (Redelegations, sdk.Error)

	QueryDelegationsTo(valAddr string) (DelegationResponses, sdk.Error)
	QueryUnbondingDelegationsFrom(valAddr string) (UnbondingDelegations, sdk.Error)
	QueryRedelegationsFrom(valAddr string) (Redelegations, sdk.Error)

	QueryValidator(address string) (Validator, sdk.Error)
	QueryValidators(page, size int) (Validators, sdk.Error)

	QueryPool() (StakePool, sdk.Error)
	QueryParams() (StakeParams, sdk.Error)
}

type StakingSubscriber interface {
	SubscribeValidatorInfoUpdates(validator string,
		handler func(data EventDataMsgEditValidator)) (sdk.Subscription, sdk.Error)
}

type Staking interface {
	sdk.Module
	StakingTx
	StakingQueries
	StakingSubscriber
}

type DelegationResponses []DelegationResponse
type DelegationResponse struct {
	Delegation Delegation `json:"delegation"`
	Balance    sdk.Coin   `json:"balance"`
}

type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
}

type UnbondingDelegations []UnbondingDelegation
type UnbondingDelegation struct {
	DelegatorAddress string                     `json:"delegator_address"`
	ValidatorAddress string                     `json:"validator_address"`
	Entries          []UnbondingDelegationEntry `json:"entries"`
}

type UnbondingDelegationEntry struct {
	CreationHeight int64     `json:"creation_height,omitempty"`
	CompletionTime time.Time `json:"completion_time"`
	InitialBalance sdk.Int   `json:"initial_balance"`
	Balance        sdk.Int   `json:"balance"`
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
	OperatorAddress   string      `json:"operator_address"`
	ConsensusPubkey   string      `json:"consensus_pubkey"`
	Jailed            bool        `json:"jailed"`
	Status            string      `json:"status"`
	Tokens            string      `json:"tokens"`
	DelegatorShares   string      `json:"delegator_shares"`
	Description       Description `json:"description"`
	UnbondingHeight   int64       `json:"unbonding_height"`
	UnbondingTime     string      `json:"unbonding_time"`
	Commission        Commission  `json:"commission"`
	MinSelfDelegation string      `json:"min_self_delegation"`
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
	CommissionRates `json:"commission_rates"`
	UpdateTime      string `json:"update_time"`
}

type CommissionRates struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
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
