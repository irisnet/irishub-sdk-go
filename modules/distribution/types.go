package distribution

import (
	"errors"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/json"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "distr"
)

var (
	_ sdk.Msg = MsgSetWithdrawAddress{}
	_ sdk.Msg = MsgWithdrawDelegatorReward{}
	_ sdk.Msg = MsgWithdrawDelegatorRewardsAll{}
	_ sdk.Msg = MsgWithdrawValidatorRewardsAll{}

	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

//______________________________________________________________________

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgSetWithdrawAddress struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	WithdrawAddr  sdk.AccAddress `json:"withdraw_addr"`
}

func (msg MsgSetWithdrawAddress) Type() string { return "set_withdraw_address" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.DelegatorAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return errors.New("delegator address is nil")
	}
	if msg.WithdrawAddr == nil {
		return errors.New("withdraw address is nil")
	}
	return nil
}

//______________________________________________________________________

// msg struct for delegation withdraw for all of the delegator's delegations
type MsgWithdrawDelegatorRewardsAll struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
}

func NewMsgWithdrawDelegatorRewardsAll(delAddr sdk.AccAddress) MsgWithdrawDelegatorRewardsAll {
	return MsgWithdrawDelegatorRewardsAll{
		DelegatorAddr: delAddr,
	}
}

func (msg MsgWithdrawDelegatorRewardsAll) Type() string { return "withdraw_delegation_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorRewardsAll) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorRewardsAll) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgWithdrawDelegatorRewardsAll) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return errors.New("delegator address is nil")
	}
	return nil
}

//______________________________________________________________________

// msg struct for delegation withdraw from a single validator
type MsgWithdrawDelegatorReward struct {
	DelegatorAddr sdk.AccAddress `json:"delegator_addr"`
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
}

func (msg MsgWithdrawDelegatorReward) Type() string { return "withdraw_delegation_reward" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgWithdrawDelegatorReward) ValidateBasic() error {
	if msg.DelegatorAddr == nil {
		return errors.New("delegator address is nil")
	}
	if msg.ValidatorAddr == nil {
		return errors.New("validator address is nil")
	}
	return nil
}

//______________________________________________________________________

// msg struct for validator withdraw
type MsgWithdrawValidatorRewardsAll struct {
	ValidatorAddr sdk.ValAddress `json:"validator_addr"`
}

func (msg MsgWithdrawValidatorRewardsAll) Type() string { return "withdraw_validator_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawValidatorRewardsAll) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr.Bytes())}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawValidatorRewardsAll) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgWithdrawValidatorRewardsAll) ValidateBasic() error {
	if msg.ValidatorAddr == nil {
		return errors.New("validator address is nil")
	}
	return nil
}

type rewards struct {
	Total       sdk.Coins            `json:"total"`
	Delegations []delegationsRewards `json:"delegations"`
	Commission  sdk.Coins            `json:"commission"`
}

func (r rewards) Convert() interface{} {
	var delegations []rpc.DelegationRewards
	for _, d := range r.Delegations {
		delegations = append(delegations, rpc.DelegationRewards{
			Validator: d.Validator.String(),
			Reward:    d.Reward,
		})
	}
	return rpc.Rewards{
		Total:       r.Total,
		Commission:  r.Commission,
		Delegations: delegations,
	}
}

type delegationsRewards struct {
	Validator sdk.ValAddress `json:"validator"`
	Reward    sdk.Coins      `json:"reward"`
}

func registerCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgWithdrawDelegatorRewardsAll{}, "irishub/distr/MsgWithdrawDelegationRewardsAll")
	cdc.RegisterConcrete(MsgWithdrawDelegatorReward{}, "irishub/distr/MsgWithdrawDelegationReward")
	cdc.RegisterConcrete(MsgWithdrawValidatorRewardsAll{}, "irishub/distr/MsgWithdrawValidatorRewardsAll")
	cdc.RegisterConcrete(MsgSetWithdrawAddress{}, "irishub/distr/MsgModifyWithdrawAddress")

	//cdc.RegisterConcrete(DelegationDistInfo{}, "irishub/distr/DelegationDistInfo")
	//cdc.RegisterConcrete(FeePool{}, "irishub/distr/FeePool")
	//cdc.RegisterConcrete(&Params{}, "irishub/distr/Params")
}
