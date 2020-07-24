package staking

import (
	"bytes"
	"errors"
	"time"

	"github.com/tendermint/tendermint/crypto"

	"github.com/irisnet/irishub-sdk-go/rpc"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	ModuleName = "stake"
)

var (
	_ sdk.Msg = MsgDelegate{}
	_ sdk.Msg = MsgUndelegate{}
	_ sdk.Msg = MsgBeginRedelegate{}

	cdc         = sdk.NewAminoCodec()
	defaultPage = 1
)

func init() {
	registerCodec(cdc)
}

// MsgCreateValidator - struct for bonding transactions
type MsgCreateValidator struct {
	Description
	Commission    CommissionMsg
	DelegatorAddr sdk.AccAddress `json:"delegator_address"`
	ValidatorAddr sdk.ValAddress `json:"validator_address"`
	PubKey        crypto.PubKey  `json:"pubkey"`
	Delegation    sdk.Coin       `json:"delegation"`
}

//nolint
func (msg MsgCreateValidator) Route() string { return ModuleName }
func (msg MsgCreateValidator) Type() string  { return "create_validator" }

// Return address(es) that must sign over msg.GetSignBytes()
func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	addrs := []sdk.AccAddress{msg.DelegatorAddr}

	if !bytes.Equal(msg.DelegatorAddr.Bytes(), msg.ValidatorAddr.Bytes()) {
		// if validator addr is not same as delegator addr, validator must sign
		// msg as well
		addrs = append(addrs, sdk.AccAddress(msg.ValidatorAddr))
	}
	return addrs
}

// get the bytes for the message signer to sign on
func (msg MsgCreateValidator) GetSignBytes() []byte {
	pukKey, err := sdk.Bech32ifyConsPub(msg.PubKey)
	if err != nil {
		panic(err)
	}

	b, err := cdc.MarshalJSON(struct {
		Description
		Commission    CommissionMsg
		DelegatorAddr sdk.AccAddress `json:"delegator_address"`
		ValidatorAddr sdk.ValAddress `json:"validator_address"`
		PubKey        string         `json:"pubkey"`
		Delegation    sdk.Coin       `json:"delegation"`
	}{
		Description:   msg.Description,
		ValidatorAddr: msg.ValidatorAddr,
		PubKey:        pukKey,
		Delegation:    msg.Delegation,
	})
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgCreateValidator) ValidateBasic() error {
	return nil
}

//______________________________________________________________________

// MsgDelegate - struct for bonding transactions
type MsgDelegate struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
	Delegation    sdk.Coin       `json:"delegation"`
}

func (msg MsgDelegate) Route() string { return ModuleName }

//nolint
func (msg MsgDelegate) Type() string { return "delegate" }
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgDelegate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
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
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
	SharesAmount  sdk.Dec        `json:"shares_amount"`
}

//nolint
func (msg MsgUndelegate) Route() string                { return ModuleName }
func (msg MsgUndelegate) Type() string                 { return "begin_unbonding" }
func (msg MsgUndelegate) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.DelegatorAddr} }

// get the bytes for the message signer to sign on
func (msg MsgUndelegate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(struct {
		DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
		ValidatorAddr sdk.ValAddress `json:"validator_addr"`
		SharesAmount  string         `json:"shares_amount"`
	}{
		DelegatorAddr: msg.DelegatorAddr,
		ValidatorAddr: msg.ValidatorAddr,
		SharesAmount:  msg.SharesAmount.String(),
	})
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgUndelegate) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return errors.New("delegator address is nil")
	}
	if msg.ValidatorAddr == nil {
		return errors.New("validator address is nil")
	}
	if msg.SharesAmount.LTE(sdk.ZeroDec()) {
		return errors.New("shares must be > 0")
	}
	return nil
}

//______________________________________________________________________
// MsgBeginRedelegate - struct for bonding transactions
type MsgBeginRedelegate struct {
	DelegatorAddr    sdk.AccAddress `json:"delegator_addr"`
	ValidatorSrcAddr sdk.ValAddress `json:"validator_src_addr"`
	ValidatorDstAddr sdk.ValAddress `json:"validator_dst_addr"`
	SharesAmount     sdk.Dec        `json:"shares_amount"`
}

func NewMsgBeginRedelegate(delAddr sdk.AccAddress, valSrcAddr,
	valDstAddr sdk.ValAddress, sharesAmount sdk.Dec) MsgBeginRedelegate {

	return MsgBeginRedelegate{
		DelegatorAddr:    delAddr,
		ValidatorSrcAddr: valSrcAddr,
		ValidatorDstAddr: valDstAddr,
		SharesAmount:     sharesAmount,
	}
}

//nolint
func (msg MsgBeginRedelegate) Route() string { return ModuleName }
func (msg MsgBeginRedelegate) Type() string  { return "begin_redelegate" }
func (msg MsgBeginRedelegate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(struct {
		DelegatorAddr    sdk.AccAddress `json:"delegator_addr"`
		ValidatorSrcAddr sdk.ValAddress `json:"validator_src_addr"`
		ValidatorDstAddr sdk.ValAddress `json:"validator_dst_addr"`
		SharesAmount     string         `json:"shares"`
	}{
		DelegatorAddr:    msg.DelegatorAddr,
		ValidatorSrcAddr: msg.ValidatorSrcAddr,
		ValidatorDstAddr: msg.ValidatorDstAddr,
		SharesAmount:     msg.SharesAmount.String(),
	})
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
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
	if msg.SharesAmount.LTE(sdk.ZeroDec()) {
		return errors.New("shares must be > 0")
	}
	return nil
}

//______________________________________________________________________

// MsgEditValidator - struct for editing a validator
type MsgEditValidator struct {
	Description
	ValidatorAddr sdk.ValAddress `json:"address"`

	// We pass a reference to the new commission rate as it's not mandatory to
	// update. If not updated, the deserialized rate will be zero with no way to
	// distinguish if an update was intended.
	//
	// REF: #2373
	CommissionRate *sdk.Dec `json:"commission_rate"`
}

//nolint
func (msg MsgEditValidator) Route() string { return ModuleName }
func (msg MsgEditValidator) Type() string  { return "edit_validator" }
func (msg MsgEditValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgEditValidator) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(struct {
		Description
		ValidatorAddr sdk.ValAddress `json:"address"`
	}{
		Description:   msg.Description,
		ValidatorAddr: msg.ValidatorAddr,
	})
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgEditValidator) ValidateBasic() error {
	if msg.ValidatorAddr == nil {
		return errors.New("nil validator address")
	}

	if msg.Description == (Description{}) {
		return errors.New("transaction must include some information to modify")
	}
	return nil
}

//===============================for query===============================
// DelegationResponse is equivalent to Delegation except that it contains a balance
// in addition to shares which is more suitable for client responses.
type delegationResponses []delegationResponse
type delegationResponse struct {
	Delegation delegation `json:"delegation"`
	Balance    sdk.Coin   `json:"balance"`
}

type delegation struct {
	DelegatorAddress sdk.AccAddress `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress `json:"validator_address"`
	Shares           sdk.Dec        `json:"shares"`
}

func (ds delegationResponses) Convert() interface{} {
	dses := make(rpc.DelegationResponses, len(ds))
	for i, v := range ds {
		dses[i] = v.Convert().(rpc.DelegationResponse)
	}
	return dses
}

func (d delegationResponse) Convert() interface{} {
	return rpc.DelegationResponse{
		Delegation: rpc.Delegation{
			DelegatorAddress: d.Delegation.DelegatorAddress.String(),
			ValidatorAddress: d.Delegation.ValidatorAddress.String(),
			Shares:           d.Delegation.Shares.String(),
		},
		Balance: d.Balance,
	}
}

// UnbondingDelegation stores all of a single delegator's unbonding bonds
// for a single validator in an time-ordered list
type unbondingDelegations []unbondingDelegation
type unbondingDelegation struct {
	DelegatorAddress sdk.AccAddress             `json:"delegator_address"`
	ValidatorAddress sdk.ValAddress             `json:"validator_address"`
	Entries          []unbondingDelegationEntry `protobuf:"bytes,3,rep,name=entries,proto3" json:"entries"`
}

// UnbondingDelegationEntry defines an unbonding object with relevant metadata.
type unbondingDelegationEntry struct {
	CreationHeight int64     `json:"creation_height,omitempty"`
	CompletionTime time.Time `json:"completion_time"`
	InitialBalance sdk.Int   `json:"initial_balance"`
	Balance        sdk.Int   `json:"balance"`
}

func (ubds unbondingDelegations) Convert() interface{} {
	uds := make(rpc.UnbondingDelegations, len(ubds))
	for i, v := range ubds {
		uds[i] = v.Convert().(rpc.UnbondingDelegation)
	}
	return uds
}

func (ubd unbondingDelegation) Convert() interface{} {
	entries := make([]rpc.UnbondingDelegationEntry, len(ubd.Entries))
	for i, v := range ubd.Entries {
		entries[i] = v.Convert().(rpc.UnbondingDelegationEntry)
	}

	return rpc.UnbondingDelegation{
		DelegatorAddress: ubd.DelegatorAddress.String(),
		ValidatorAddress: ubd.ValidatorAddress.String(),
		Entries:          entries,
	}
}

func (ubdEntry unbondingDelegationEntry) Convert() interface{} {
	return rpc.UnbondingDelegationEntry{
		CreationHeight: ubdEntry.CreationHeight,
		CompletionTime: ubdEntry.CompletionTime,
		InitialBalance: ubdEntry.InitialBalance,
		Balance:        ubdEntry.Balance,
	}
}

// redelegation reflects a delegation's passive re-delegation queue.
type redelegation struct {
	DelegatorAddr    sdk.AccAddress `json:"delegator_addr"`     // delegator
	ValidatorSrcAddr sdk.ValAddress `json:"validator_src_addr"` // validator redelegation source operator addr
	ValidatorDstAddr sdk.ValAddress `json:"validator_dst_addr"` // validator redelegation destination operator addr
	CreationHeight   int64          `json:"creation_height"`    // height which the redelegation took place
	MinTime          time.Time      `json:"min_time"`           // unix time for redelegation completion
	InitialBalance   sdk.Coin       `json:"initial_balance"`    // initial balance when redelegation started
	Balance          sdk.Coin       `json:"balance"`            // current balance
	SharesSrc        sdk.Dec        `json:"shares_src"`         // amount of source shares redelegating
	SharesDst        sdk.Dec        `json:"shares_dst"`         // amount of destination shares redelegating
}

func (d redelegation) Convert() interface{} {
	return rpc.Redelegation{
		DelegatorAddr:    d.DelegatorAddr.String(),
		ValidatorSrcAddr: d.ValidatorDstAddr.String(),
		ValidatorDstAddr: d.ValidatorDstAddr.String(),
		CreationHeight:   d.CreationHeight,
		MinTime:          d.MinTime.String(),
		InitialBalance:   d.InitialBalance,
		Balance:          d.Balance,
		SharesSrc:        d.SharesSrc.String(),
		SharesDst:        d.SharesDst.String(),
	}
}

type redelegations []redelegation

func (ds redelegations) Convert() interface{} {
	rds := make(rpc.Redelegations, len(ds))
	for _, d := range ds {
		rds = append(rds, d.Convert().(rpc.Redelegation))
	}
	return rds
}

type validator struct {
	OperatorAddress   string      `json:"operator_address"` // address of the validator's operator; bech encoded in JSON
	ConsensusPubkey   string      `json:"consensus_pubkey"` // the consensus public key of the validator; bech encoded in JSON
	Jailed            bool        `json:"jailed"`           // has the validator been jailed from bonded status?
	Status            bondStatus  `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Tokens            string      `json:"tokens"`           // delegated tokens (incl. self-delegation)
	Description       Description `json:"description"`      // description terms for the validator
	UnbondingHeight   int64       `json:"unbonding_height"` // if unbonding, height at which this validator has begun unbonding
	UnbondingTime     time.Time   `json:"unbonding_time"`   // if unbonding, min time for the validator to complete unbonding
	Commission        Commission  `json:"commission"`       // commission parameters
	MinSelfDelegation string      `json:"min_self_delegation"`
}

type validators []validator

func (v validator) Convert() interface{} {
	return rpc.Validator{
		OperatorAddress: v.OperatorAddress,
		ConsensusPubkey: v.ConsensusPubkey,
		Jailed:          v.Jailed,
		Status:          v.Status.String(),
		Tokens:          v.Tokens,
		Description: rpc.Description{
			Moniker:  v.Description.Moniker,
			Identity: v.Description.Identity,
			Website:  v.Description.Website,
			Details:  v.Description.Details,
		},
		UnbondingHeight: v.UnbondingHeight,
		UnbondingTime:   v.UnbondingTime.String(),
		Commission: rpc.Commission{
			CommissionRates: rpc.CommissionRates{
				Rate:          v.Commission.Rate.String(),
				MaxRate:       v.Commission.MaxRate.String(),
				MaxChangeRate: v.Commission.MaxChangeRate.String(),
			},
			UpdateTime: v.Commission.UpdateTime.String(),
		},
		MinSelfDelegation: v.MinSelfDelegation,
	}
}

func (vs validators) Convert() interface{} {
	validators := make(rpc.Validators, len(vs))
	for i, v := range vs {
		validators[i] = v.Convert().(rpc.Validator)
	}
	return validators
}

// status of a validator
type bondStatus int32

// nolint
const (
	Unbonded  bondStatus = 1
	Unbonding bondStatus = 2
	Bonded    bondStatus = 3

	BondStatusUnbonded  = "Unbonded"
	BondStatusUnbonding = "Unbonding"
	BondStatusBonded    = "Bonded"
)

func (b bondStatus) String() string {
	switch b {
	case Unbonded:
		return "Unbonded"
	case Unbonding:
		return "Unbonding"
	case Bonded:
		return "Bonded"
	default:
		panic("improper use of BondStatusToString")
	}
}

// Description - description fields for a validator
type Description struct {
	Moniker         string `json:"moniker"`          // name
	Identity        string `json:"identity"`         // optional identity signature (ex. UPort or Keybase)
	Website         string `json:"website"`          // optional website link
	SecurityContact string `json:"security_contact"` // optional securityContact
	Details         string `json:"details"`          // optional details
}

// Commission defines a commission parameters for a given validator.
type Commission struct {
	Rate          sdk.Dec   `json:"rate"`            // the commission rate charged to delegators
	MaxRate       sdk.Dec   `json:"max_rate"`        // maximum commission rate which validator can ever charge
	MaxChangeRate sdk.Dec   `json:"max_change_rate"` // maximum daily increase of the validator commission
	UpdateTime    time.Time `json:"update_time"`     // the last time the commission rate was changed
}

// CommissionMsg defines a commission message to be used for creating a
// validator.
type CommissionMsg struct {
	Rate          sdk.Dec `json:"rate"`            // the commission rate charged to delegators
	MaxRate       sdk.Dec `json:"max_rate"`        // maximum commission rate which validator can ever charge
	MaxChangeRate sdk.Dec `json:"max_change_rate"` // maximum daily increase of the validator commission
}

type Pool struct {
	LooseTokens  sdk.Dec `json:"loose_tokens"`  // tokens which are not bonded in a validator
	BondedTokens sdk.Dec `json:"bonded_tokens"` // reserve of bonded tokens
}

func (p Pool) Convert() interface{} {
	return rpc.StakePool{
		LooseTokens:  p.LooseTokens.String(),
		BondedTokens: p.BondedTokens.String(),
	}
}

// params defines the high level settings for staking
type params struct {
	UnbondingTime time.Duration `json:"unbonding_time"`
	MaxValidators uint16        `json:"max_validators"` // maximum number of validators
}

func (p params) Convert() interface{} {
	return rpc.StakeParams{
		UnbondingTime: p.UnbondingTime.String(),
		MaxValidators: int(p.MaxValidators),
	}
}

func registerCodec(cdc sdk.Codec) {
	//cdc.RegisterConcrete(Pool{}, "irishub/stake/Pool")
	cdc.RegisterConcrete(&params{}, "irishub/stake/Params")
	cdc.RegisterConcrete(validator{}, "custom/staking/Validator")
	cdc.RegisterConcrete(delegation{}, "irishub/stake/Delegation")
	cdc.RegisterConcrete(unbondingDelegation{}, "irishub/stake/UnbondingDelegation")
	cdc.RegisterConcrete(redelegation{}, "irishub/stake/Redelegation")

	cdc.RegisterConcrete(MsgCreateValidator{}, "irishub/stake/MsgCreateValidator")
	cdc.RegisterConcrete(MsgEditValidator{}, "irishub/stake/MsgEditValidator")
	cdc.RegisterConcrete(MsgDelegate{}, "irishub/stake/MsgDelegate")
	cdc.RegisterConcrete(MsgUndelegate{}, "irishub/stake/BeginUnbonding") //TODO
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "irishub/stake/BeginRedelegate")
}
