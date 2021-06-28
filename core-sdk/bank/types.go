package bank

import (
	"errors"
	"fmt"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	maxMsgLen  = 5
	ModuleName = "bank"

	TypeMsgSend      = "send"
	TypeMsgMultiSend = "multisend"
)

var _ sdk.Msg = &MsgSend{}

// NewMsgSend - construct a msg to send coins from one account to another.
//nolint:interfacer
func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) *MsgSend {
	return &MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      amount,
	}
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgSend) Type() string { return TypeMsgSend }

func (msg MsgSend) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return errors.New("invalid sender address")
	}

	_, err = sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return errors.New("invalid recipient address")
	}

	if !msg.Amount.IsValid() {
		return errors.New("invalid coins")
	}

	if !msg.Amount.IsAllPositive() {
		return errors.New("invalid coins")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

var _ sdk.Msg = &MsgMultiSend{}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) *MsgMultiSend {
	return &MsgMultiSend{Inputs: in, Outputs: out}
}

func (msg MsgMultiSend) Route() string { return ModuleName }

// Type Implements Msg
func (msg MsgMultiSend) Type() string { return TypeMsgMultiSend }

// Implements Msg.
func (msg MsgMultiSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return errors.New("invalid input coins")
	}
	if len(msg.Outputs) == 0 {
		return errors.New("invalid output coins")
	}
	// make sure all inputs and outputs are individually valid
	var totalIn, totalOut sdk.Coins
	for _, in := range msg.Inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}
		totalIn = totalIn.Add(in.Coins...)
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}
		totalOut = totalOut.Add(out.Coins...)
	}
	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return errors.New("inputs and outputs don't match")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners Implements Msg.
func (msg MsgMultiSend) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addr, _ := sdk.AccAddressFromBech32(in.Address)
		addrs[i] = addr
	}

	return addrs
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(in.Address)
	if err != nil {
		return err
	}

	if in.Coins.Empty() {
		return errors.New("empty input coins")
	}

	if !in.Coins.IsValid() {
		return fmt.Errorf("invalid input coins [%s]", in.Coins)
	}

	if !in.Coins.IsAllPositive() {
		return fmt.Errorf("invalid input coins [%s]", in.Coins)
	}

	return nil
}

// NewInput - create a transaction input, used with MsgMultiSend
//nolint:interfacer
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	return Input{
		Address: addr.String(),
		Coins:   coins,
	}
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(out.Address)
	if err != nil {
		return fmt.Errorf("invalid output address (%s)", err)
	}

	if out.Coins.Empty() {
		return errors.New("empty input coins")
	}

	if !out.Coins.IsValid() {
		return fmt.Errorf("invalid input coins [%s]", out.Coins)
	}

	if !out.Coins.IsAllPositive() {
		return fmt.Errorf("invalid input coins [%s]", out.Coins)
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgMultiSend
//nolint:interfacer
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	return Output{
		Address: addr.String(),
		Coins:   coins,
	}
}

// ValidateInputsOutputs validates that each respective input and output is
// valid and that the sum of inputs is equal to the sum of outputs.
func ValidateInputsOutputs(inputs []Input, outputs []Output) error {
	var totalIn, totalOut sdk.Coins
	for _, in := range inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}

		totalIn = totalIn.Add(in.Coins...)
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}

		totalOut = totalOut.Add(out.Coins...)
	}

	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return errors.New("sum inputs != sum outputs")
	}

	return nil
}
