package types

import (
	"errors"
	"time"

	"github.com/irisnet/irishub-sdk-go/utils"

	"github.com/tendermint/go-amino"
)

var (
	_ Msg = MsgDelegate{}
	_ Msg = MsgUndelegate{}
	_ Msg = MsgBeginRedelegate{}
)

//______________________________________________________________________

// MsgDelegate - struct for bonding transactions
type MsgDelegate struct {
	DelegatorAddr AccAddress `json:"delegator_addr"`
	ValidatorAddr ValAddress `json:"validator_addr"`
	Delegation    Coin       `json:"delegation"`
}

func NewMsgDelegate(delAddr AccAddress, valAddr ValAddress, delegation Coin) MsgDelegate {
	return MsgDelegate{
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
		Delegation:    delegation,
	}
}

//nolint
func (msg MsgDelegate) Type() string { return "delegate" }
func (msg MsgDelegate) GetSigners() []AccAddress {
	return []AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgDelegate) GetSignBytes() []byte {
	b, err := defaultCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// quick validity check
func (msg MsgDelegate) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return errors.New("delegator address is nil")
	}
	if msg.ValidatorAddr == nil {
		return errors.New("validator address is nil")
	}
	if !msg.Delegation.IsValidIrisAtto() {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

//______________________________________________________________________

// MsgUndelegate - struct for unbonding transactions
type MsgUndelegate struct {
	DelegatorAddr AccAddress `json:"delegator_addr"`
	ValidatorAddr ValAddress `json:"validator_addr"`
	SharesAmount  Dec        `json:"shares_amount"`
}

func NewMsgUndelegate(delAddr AccAddress, valAddr ValAddress, sharesAmount Dec) MsgUndelegate {
	return MsgUndelegate{
		DelegatorAddr: delAddr,
		ValidatorAddr: valAddr,
		SharesAmount:  sharesAmount,
	}
}

//nolint
func (msg MsgUndelegate) Type() string             { return "begin_unbonding" }
func (msg MsgUndelegate) GetSigners() []AccAddress { return []AccAddress{msg.DelegatorAddr} }

// get the bytes for the message signer to sign on
func (msg MsgUndelegate) GetSignBytes() []byte {
	b, err := defaultCdc.MarshalJSON(struct {
		DelegatorAddr AccAddress `json:"delegator_addr"`
		ValidatorAddr ValAddress `json:"validator_addr"`
		SharesAmount  string     `json:"shares_amount"`
	}{
		DelegatorAddr: msg.DelegatorAddr,
		ValidatorAddr: msg.ValidatorAddr,
		SharesAmount:  msg.SharesAmount.String(),
	})
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// quick validity check
func (msg MsgUndelegate) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return errors.New("delegator address is nil")
	}
	if msg.ValidatorAddr == nil {
		return errors.New("validator address is nil")
	}
	if msg.SharesAmount.Int == nil || msg.SharesAmount.LTE(ZeroDec()) {
		return errors.New("shares must be > 0")
	}
	return nil
}

// MsgBeginRedelegate - struct for bonding transactions
type MsgBeginRedelegate struct {
	DelegatorAddr    AccAddress `json:"delegator_addr"`
	ValidatorSrcAddr ValAddress `json:"validator_src_addr"`
	ValidatorDstAddr ValAddress `json:"validator_dst_addr"`
	SharesAmount     Dec        `json:"shares_amount"`
}

func NewMsgBeginRedelegate(delAddr AccAddress, valSrcAddr,
	valDstAddr ValAddress, sharesAmount Dec) MsgBeginRedelegate {

	return MsgBeginRedelegate{
		DelegatorAddr:    delAddr,
		ValidatorSrcAddr: valSrcAddr,
		ValidatorDstAddr: valDstAddr,
		SharesAmount:     sharesAmount,
	}
}

//nolint
func (msg MsgBeginRedelegate) Type() string { return "begin_redelegate" }
func (msg MsgBeginRedelegate) GetSigners() []AccAddress {
	return []AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	b, err := defaultCdc.MarshalJSON(struct {
		DelegatorAddr    AccAddress `json:"delegator_addr"`
		ValidatorSrcAddr ValAddress `json:"validator_src_addr"`
		ValidatorDstAddr ValAddress `json:"validator_dst_addr"`
		SharesAmount     string     `json:"shares"`
	}{
		DelegatorAddr:    msg.DelegatorAddr,
		ValidatorSrcAddr: msg.ValidatorSrcAddr,
		ValidatorDstAddr: msg.ValidatorDstAddr,
		SharesAmount:     msg.SharesAmount.String(),
	})
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// quick validity check
func (msg MsgBeginRedelegate) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return errors.New("delegator address is nil")
	}
	if msg.ValidatorSrcAddr == nil {
		return errors.New("validator address is nil")
	}
	if msg.ValidatorDstAddr == nil {
		return errors.New("validator address is nil")
	}
	if msg.SharesAmount.Int == nil || msg.SharesAmount.LTE(ZeroDec()) {
		return errors.New("shares must be > 0")
	}
	return nil
}

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
	Tokens          Dec         `json:"tokens"`
	DelegatorShares Dec         `json:"delegator_shares"`
	Description     Description `json:"description"`
	BondHeight      int64       `json:"bond_height"`
	UnbondingHeight int64       `json:"unbonding_height"`
	UnbondingTime   time.Time   `json:"unbonding_time"`
	Commission      Commission  `json:"commission"`
}

// DelegatorShareExRate gets the exchange rate of tokens over delegator shares.
// UNITS: tokens/delegator-shares
func (v Validator) DelegatorShareExRate() Dec {
	if v.DelegatorShares.IsZero() {
		return OneDec()
	}
	return v.Tokens.Quo(v.DelegatorShares)
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

	cdc.RegisterConcrete(MsgDelegate{}, "irishub/stake/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "irishub/stake/BeginUnbonding", nil) //TODO
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "irishub/stake/BeginRedelegate", nil)
}
