package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/irisnet/irishub-sdk-go/utils"
	"github.com/tendermint/go-amino"
)

const memoRegexpLengthLimit = 50

var (
	_ Msg = MsgSend{}
	_ Msg = MsgBurn{}
	_ Msg = MsgSetMemoRegexp{}
)

type MsgSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) MsgSend {
	return MsgSend{Inputs: in, Outputs: out}
}

// Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// Implements Msg.
func (msg MsgSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return errors.New("invalid input coins")
	}
	if len(msg.Outputs) == 0 {
		return errors.New("invalid output coins")
	}
	// make sure all inputs and outputs are individually valid
	var totalIn, totalOut Coins
	for _, in := range msg.Inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}
		totalIn = totalIn.Add(in.Coins)
	}
	for _, out := range msg.Outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}
		totalOut = totalOut.Add(out.Coins)
	}
	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return errors.New("inputs and outputs don't match")
	}
	return nil
}

// Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	var inputs, outputs []json.RawMessage
	for _, input := range msg.Inputs {
		inputs = append(inputs, input.GetSignBytes())
	}
	for _, output := range msg.Outputs {
		outputs = append(outputs, output.GetSignBytes())
	}
	b, err := DefaultCodec().MarshalJSON(struct {
		Inputs  []json.RawMessage `json:"inputs"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Inputs:  inputs,
		Outputs: outputs,
	})
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSend) GetSigners() []AccAddress {
	addrs := make([]AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

//----------------------------------------
// Input

// Transaction Input
type Input struct {
	Address AccAddress `json:"address"`
	Coins   Coins      `json:"coins"`
}

// Return bytes to sign for Input
func (in Input) GetSignBytes() []byte {
	bin, err := DefaultCodec().MarshalJSON(in)
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(bin)
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if len(in.Address) == 0 {
		return errors.New(fmt.Sprintf(fmt.Sprintf("account %s is invalid", in.Address.String())))
	}
	if in.Coins.Empty() {
		return errors.New("empty input coins")
	}
	if !in.Coins.IsValid() {
		return errors.New(fmt.Sprintf("invalid input coins [%s]", in.Coins))
	}
	return nil
}

// NewInput - create a transaction input, used with MsgSend
func NewInput(addr AccAddress, coins Coins) Input {
	input := Input{
		Address: addr,
		Coins:   coins,
	}
	return input
}

//----------------------------------------
// Output

// Transaction Output
type Output struct {
	Address AccAddress `json:"address"`
	Coins   Coins      `json:"coins"`
}

// Return bytes to sign for Output
func (out Output) GetSignBytes() []byte {
	bin, err := DefaultCodec().MarshalJSON(out)
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(bin)
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	if len(out.Address) == 0 {
		return errors.New(fmt.Sprintf(fmt.Sprintf("account %s is invalid", out.Address.String())))
	}
	if out.Coins.Empty() {
		return errors.New("empty input coins")
	}
	if !out.Coins.IsValid() {
		return errors.New(fmt.Sprintf("invalid input coins [%s]", out.Coins))
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgSend
func NewOutput(addr AccAddress, coins Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}

// MsgBurn - high level transaction of the coin module
type MsgBurn struct {
	Owner AccAddress `json:"owner"`
	Coins Coins      `json:"coins"`
}

// NewMsgBurn - construct MsgBurn
func NewMsgBurn(owner AccAddress, coins Coins) MsgBurn {
	return MsgBurn{Owner: owner, Coins: coins}
}

// Implements Msg.
// nolint
func (msg MsgBurn) Route() string { return "bank" }
func (msg MsgBurn) Type() string  { return "burn" }

// Implements Msg.
func (msg MsgBurn) ValidateBasic() error {
	if len(msg.Owner) == 0 {
		return errors.New(fmt.Sprintf("invalid address:%s", msg.Owner.String()))
	}
	if msg.Coins.Empty() {
		return errors.New("empty coins to burn")
	}
	if !msg.Coins.IsValid() {
		return errors.New(fmt.Sprintf("invalid coins to burn [%s]", msg.Coins))
	}
	return nil
}

// Implements Msg.
func (msg MsgBurn) GetSignBytes() []byte {
	b, err := defaultCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgBurn) GetSigners() []AccAddress {
	return []AccAddress{msg.Owner}
}

// MsgSetMemoRegexp - set memo regexp
type MsgSetMemoRegexp struct {
	Owner      AccAddress `json:"owner"`
	MemoRegexp string     `json:"memo_regexp"`
}

// NewMsgSetMemoRegexp - construct MsgSetMemoRegexp
func NewMsgSetMemoRegexp(owner AccAddress, memoRegexp string) MsgSetMemoRegexp {
	return MsgSetMemoRegexp{Owner: owner, MemoRegexp: memoRegexp}
}

// Implements Msg.
// nolint
func (msg MsgSetMemoRegexp) Type() string { return "set-memo-regexp" }

// Implements Msg.
func (msg MsgSetMemoRegexp) ValidateBasic() error {
	if len(msg.Owner) == 0 {
		return errors.New(fmt.Sprintf("invalid address:%s", msg.Owner.String()))
	}
	if len(msg.MemoRegexp) > memoRegexpLengthLimit {
		return errors.New("memo regexp length exceeds limit")
	}
	if _, err := regexp.Compile(msg.MemoRegexp); err != nil {
		return errors.New("invalid memo regexp")
	}
	return nil
}

// Implements Msg.
func (msg MsgSetMemoRegexp) GetSignBytes() []byte {
	b, err := defaultCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSetMemoRegexp) GetSigners() []AccAddress {
	return []AccAddress{msg.Owner}
}

type TokenStats struct {
	LooseTokens  Coins `json:"loose_tokens"`
	BondedTokens Coins `json:"bonded_tokens"`
	BurnedTokens Coins `json:"burned_tokens"`
	TotalSupply  Coins `json:"total_supply"`
}

func RegisterBank(cdc *amino.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "irishub/bank/Send", nil)
	cdc.RegisterConcrete(MsgBurn{}, "irishub/bank/Burn", nil)
	cdc.RegisterConcrete(MsgSetMemoRegexp{}, "irishub/bank/SetMemoRegexp", nil)
}
