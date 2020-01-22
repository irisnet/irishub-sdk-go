package types

import (
	"time"

	"github.com/tendermint/go-amino"
)

// status of a validator
type BondStatus byte

// nolint
const (
	Unbonded  BondStatus = 0x00
	Unbonding BondStatus = 0x01
	Bonded    BondStatus = 0x02
)

//BondStatusToString for pretty prints of Bond Status
func BondStatusToString(b BondStatus) string {
	switch b {
	case 0x00:
		return "Unbonded"
	case 0x01:
		return "Unbonding"
	case 0x02:
		return "Bonded"
	default:
		panic("improper use of BondStatusToString")
	}
}

type Delegations []Delegation
type Delegation struct {
	DelegatorAddr string `json:"delegator_addr"`
	ValidatorAddr string `json:"validator_addr"`
	Shares        string `json:"shares"`
	Height        int64  `json:"height"`
}

type UnbondingDelegations []UnbondingDelegation
type UnbondingDelegation struct {
	TxHash         string    `json:"tx_hash"`
	DelegatorAddr  string    `json:"delegator_addr"`
	ValidatorAddr  string    `json:"validator_addr"`
	CreationHeight int64     `json:"creation_height"`
	MinTime        time.Time `json:"min_time"`
	InitialBalance Coin      `json:"initial_balance"`
	Balance        Coin      `json:"balance"`
}

type Redelegations []Redelegation
type Redelegation struct {
	DelegatorAddr    string `json:"delegator_addr"`
	ValidatorSrcAddr string `json:"validator_src_addr"`
	ValidatorDstAddr string `json:"validator_dst_addr"`
	CreationHeight   int64  `json:"creation_height"`
	MinTime          string `json:"min_time"`
	InitialBalance   Coin   `json:"initial_balance"`
	Balance          Coin   `json:"balance"`
	SharesSrc        string `json:"shares_src"`
	SharesDst        string `json:"shares_dst"`
}

type Validators []Validator
type Validator struct {
	OperatorAddress string      `json:"operator_address"`
	ConsensusPubkey string      `json:"consensus_pubkey"`
	Jailed          bool        `json:"jailed"`
	Status          BondStatus  `json:"status"`
	Tokens          string      `json:"tokens"`
	DelegatorShares string      `json:"delegator_shares"`
	Description     Description `json:"description"`
	BondHeight      int64       `json:"bond_height"`
	UnbondingHeight int64       `json:"unbonding_height"`
	UnbondingTime   time.Time   `json:"unbonding_time"`
	Commission      Commission  `json:"commission"`
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

func RegisterStake(cdc *amino.Codec) {
	cdc.RegisterConcrete(Validator{}, "irishub/stake/Validator", nil)
	cdc.RegisterConcrete(Delegation{}, "irishub/stake/Delegation", nil)
	cdc.RegisterConcrete(UnbondingDelegation{}, "irishub/stake/UnbondingDelegation", nil)
	cdc.RegisterConcrete(Redelegation{}, "irishub/stake/Redelegation", nil)
}
