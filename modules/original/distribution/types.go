package distribution

import (
	"errors"
	"github.com/irisnet/irishub-sdk-go/types/original"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	ModuleName = "distr"
)

var (
	_ original.Msg = MsgSetWithdrawAddress{}
	_ original.Msg = MsgWithdrawDelegatorReward{}
	_ original.Msg = MsgWithdrawDelegatorRewardsAll{}
	_ original.Msg = MsgWithdrawValidatorCommission{}

	cdc = original.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

//______________________________________________________________________

// msg struct for changing the withdraw address for a delegator (or validator self-delegation)
type MsgSetWithdrawAddress struct {
	DelegatorAddr original.AccAddress `json:"delegator_addr"`
	WithdrawAddr  original.AccAddress `json:"withdraw_addr"`
}

func (msg MsgSetWithdrawAddress) Route() string { return ModuleName }

func (msg MsgSetWithdrawAddress) Type() string { return "set_withdraw_address" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgSetWithdrawAddress) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.DelegatorAddr}
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

// msg struct for delegation withdraw for all of the delegator's delegations
type MsgWithdrawDelegatorRewardsAll struct {
	DelegatorAddr original.AccAddress `json:"delegator_addr"`
}

func NewMsgWithdrawDelegatorRewardsAll(delAddr original.AccAddress) MsgWithdrawDelegatorRewardsAll {
	return MsgWithdrawDelegatorRewardsAll{
		DelegatorAddr: delAddr,
	}
}

func (msg MsgWithdrawDelegatorRewardsAll) Route() string { return ModuleName }

func (msg MsgWithdrawDelegatorRewardsAll) Type() string { return "withdraw_delegation_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorRewardsAll) GetSigners() []original.AccAddress {
	return []original.AccAddress{original.AccAddress(msg.DelegatorAddr)}
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
	DelegatorAddr original.AccAddress `json:"delegator_addr"`
	ValidatorAddr original.ValAddress `json:"validator_addr"`
}

func (msg MsgWithdrawDelegatorReward) Route() string { return ModuleName }

func (msg MsgWithdrawDelegatorReward) Type() string { return "withdraw_delegation_reward" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDelegatorReward) GetSigners() []original.AccAddress {
	return []original.AccAddress{original.AccAddress(msg.DelegatorAddr)}
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

// msg struct for validator withdraw
type MsgWithdrawValidatorCommission struct {
	ValidatorAddr original.ValAddress `json:"validator_addr"`
}

func (msg MsgWithdrawValidatorCommission) Route() string { return ModuleName }

func (msg MsgWithdrawValidatorCommission) Type() string { return "withdraw_validator_rewards_all" }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawValidatorCommission) GetSigners() []original.AccAddress {
	return []original.AccAddress{original.AccAddress(msg.ValidatorAddr.Bytes())}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawValidatorCommission) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// quick validity check
func (msg MsgWithdrawValidatorCommission) ValidateBasic() error {
	if msg.ValidatorAddr == nil {
		return errors.New("validator address is nil")
	}
	return nil
}

type rewardsResponse struct {
	Rewards []delegationDelegatorReward `json:"rewards"`
	Total   original.DecCoins           `json:"total"`
}

func (r rewardsResponse) Convert() interface{} {
	var rewards []rpc.DelegationsRewards
	for _, d := range r.Rewards {
		rewards = append(rewards, rpc.DelegationsRewards{
			Validator: d.Validator.String(),
			Reward:    d.Reward,
		})
	}
	return rpc.Rewards{
		Total:   r.Total,
		Rewards: rewards,
	}
}

type delegationDelegatorReward struct {
	Validator original.ValAddress `json:"validator_address"`
	Reward    original.DecCoins   `json:"reward"`
}

type validatorAccumulatedCommission struct {
	Commission original.DecCoins `json:"commission"`
}

func (v validatorAccumulatedCommission) Convert() interface{} {
	return rpc.ValidatorAccumulatedCommission{
		Commission: v.Commission,
	}
}

func registerCodec(cdc original.Codec) {
	cdc.RegisterConcrete(&MsgWithdrawDelegatorRewardsAll{}, "cosmos-sdk/MsgWithdrawDelegationReward")
	cdc.RegisterConcrete(&MsgWithdrawValidatorCommission{}, "cosmos-sdk/MsgWithdrawValidatorCommission")
	cdc.RegisterConcrete(&MsgSetWithdrawAddress{}, "cosmos-sdk/MsgModifyWithdrawAddress")
}
