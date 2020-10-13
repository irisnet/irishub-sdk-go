package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
	"time"
)

type StakingTx interface {
	Delegate(valAddr string, amount original.DecCoin, baseTx original.BaseTx) (original.ResultTx, original.Error)
	Undelegate(valAddr string, amount original.DecCoin, baseTx original.BaseTx) (original.ResultTx, original.Error)
	Redelegate(srcValAddr,
		dstValAddr string,
		amount original.DecCoin,
		baseTx original.BaseTx) (original.ResultTx, original.Error)
}

type StakingQueries interface {
	QueryDelegation(delAddr, valAddr string) (Delegation, original.Error)
	QueryDelegations(delAddr string) (DelegationResponses, original.Error)
	QueryDelegationsTo(valAddr string) (DelegationResponses, original.Error)

	QueryUnbondingDelegations(delAddr string) (UnbondingDelegations, original.Error)
	QueryUnbondingDelegationsFrom(valAddr string) (UnbondingDelegations, original.Error)

	QueryRedelegationsFrom(valAddr string) (RedelegationResponses, original.Error)

	QueryValidator(address string) (Validator, original.Error)
	QueryValidators(page, size int) (Validators, original.Error)

	QueryPool() (StakePool, original.Error)
	QueryParams() (StakeParams, original.Error)
}

type StakingSubscriber interface {
	SubscribeValidatorInfoUpdates(validator string,
		handler func(data EventDataMsgEditValidator)) (original.Subscription, original.Error)
}

type Staking interface {
	original.Module
	StakingTx
	StakingQueries
	StakingSubscriber
}

type DelegationResponses []DelegationResponse
type DelegationResponse struct {
	Delegation Delegation    `json:"delegation"`
	Balance    original.Coin `json:"balance"`
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
	CreationHeight int64        `json:"creation_height,omitempty"`
	CompletionTime time.Time    `json:"completion_time"`
	InitialBalance original.Int `json:"initial_balance"`
	Balance        original.Int `json:"balance"`
}

type RedelegationResponses []RedelegationResponse
type RedelegationResponse struct {
	Redelegation Redelegation                `json:"redelegation"`
	Entries      []RedelegationEntryResponse `json:"entries"`
}

type Redelegation struct {
	DelegatorAddress    string              `json:"delegator_address"`
	ValidatorSrcAddress string              `json:"validator_src_address,omitempty"`
	ValidatorDstAddress string              `json:"validator_dst_address"`
	Entries             []RedelegationEntry `json:"entries"`
}

type RedelegationEntryResponse struct {
	RedelegationEntry RedelegationEntry `json:"redelegation_entry"`
	Balance           original.Int      `json:"balance"`
}

type RedelegationEntry struct {
	CreationHeight int32        `json:"creation_height"`
	CompletionTime time.Time    `json:"completion_time"`
	InitialBalance original.Int `json:"initial_balance"`
	SharesDst      original.Dec `json:"shares_dst"`
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
func (v Validator) DelegatorShareExRate() original.Dec {
	delegatorShares, err := original.NewDecFromStr(v.DelegatorShares)
	if err != nil {
		return original.ZeroDec()
	}

	tokens, err := original.NewDecFromStr(v.Tokens)
	if err != nil {
		return original.ZeroDec()
	}
	if delegatorShares.IsZero() {
		return original.OneDec()
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
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
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
