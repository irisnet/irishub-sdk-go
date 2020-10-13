package bank

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"regexp"

	"github.com/irisnet/irishub-sdk-go/rpc"

	json2 "github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	memoRegexpLengthLimit = 50
	maxMsgLen             = 5
	ModuleName            = "bank"
)

var (
	_ original.Msg = MsgMultiSend{}
	_ original.Msg = MsgBurn{}
	_ original.Msg = MsgSetMemoRegexp{}

	cdc = original.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	FromAddress original.AccAddress `json:"from_address"`
	ToAddress   original.AccAddress `json:"to_address"`
	Amount      original.Coins      `json:"amount"`
}

type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) MsgMultiSend {
	return MsgMultiSend{Inputs: in, Outputs: out}
}

func (msg MsgMultiSend) Route() string { return ModuleName }

// Implements Msg.
func (msg MsgMultiSend) Type() string { return "send" }

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
	var totalIn, totalOut original.Coins
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

// Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	var inputs, outputs []json.RawMessage
	for _, input := range msg.Inputs {
		inputs = append(inputs, input.GetSignBytes())
	}
	for _, output := range msg.Outputs {
		outputs = append(outputs, output.GetSignBytes())
	}
	b, err := cdc.MarshalJSON(struct {
		Inputs  []json.RawMessage `json:"inputs"`
		Outputs []json.RawMessage `json:"outputs"`
	}{
		Inputs:  inputs,
		Outputs: outputs,
	})
	if err != nil {
		panic(err)
	}
	return json2.MustSort(b)
}

// Implements Msg.
func (msg MsgMultiSend) GetSigners() []original.AccAddress {
	addrs := make([]original.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

//----------------------------------------
// Input

// Transaction Input
type Input struct {
	Address original.AccAddress `json:"address"`
	Coins   original.Coins      `json:"coins"`
}

// Return bytes to sign for Input
func (in Input) GetSignBytes() []byte {
	bin, err := cdc.MarshalJSON(in)
	if err != nil {
		panic(err)
	}
	return json2.MustSort(bin)
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
func NewInput(addr original.AccAddress, coins original.Coins) Input {
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
	Address original.AccAddress `json:"address"`
	Coins   original.Coins      `json:"coins"`
}

// Return bytes to sign for Output
func (out Output) GetSignBytes() []byte {
	bin, err := cdc.MarshalJSON(out)
	if err != nil {
		panic(err)
	}
	return json2.MustSort(bin)
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
func NewOutput(addr original.AccAddress, coins original.Coins) Output {
	output := Output{
		Address: addr,
		Coins:   coins,
	}
	return output
}

// MsgBurn - high level transaction of the coin module
type MsgBurn struct {
	Owner original.AccAddress `json:"owner"`
	Coins original.Coins      `json:"coins"`
}

// NewMsgBurn - construct MsgBurn
func NewMsgBurn(owner original.AccAddress, coins original.Coins) MsgBurn {
	return MsgBurn{Owner: owner, Coins: coins}
}

// Implements Msg.
// nolint
func (msg MsgBurn) Route() string { return ModuleName }
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
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json2.MustSort(b)
}

// Implements Msg.
func (msg MsgBurn) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Owner}
}

// MsgSetMemoRegexp - set memo regexp
type MsgSetMemoRegexp struct {
	Owner      original.AccAddress `json:"owner"`
	MemoRegexp string              `json:"memo_regexp"`
}

// NewMsgSetMemoRegexp - construct MsgSetMemoRegexp
func NewMsgSetMemoRegexp(owner original.AccAddress, memoRegexp string) MsgSetMemoRegexp {
	return MsgSetMemoRegexp{Owner: owner, MemoRegexp: memoRegexp}
}

func (msg MsgSetMemoRegexp) Route() string { return ModuleName }

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
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json2.MustSort(b)
}

// Implements Msg.
func (msg MsgSetMemoRegexp) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Owner}
}

// params defines the high level settings for auth
type Params struct {
	GasPriceThreshold original.Int `json:"gas_price_threshold"` // gas price threshold
	TxSizeLimit       uint64       `json:"tx_size"`             // tx size limit
}

func (p Params) Convert() interface{} {
	return p
}

type tokenStats struct {
	LooseTokens  original.Coins `json:"loose_tokens"`
	BondedTokens original.Coins `json:"bonded_tokens"`
	BurnedTokens original.Coins `json:"burned_tokens"`
	TotalSupply  original.Coins `json:"total_supply"`
}

func (ts tokenStats) Convert() interface{} {
	return rpc.TokenStats{
		LooseTokens:  ts.LooseTokens,
		BondedTokens: ts.BondedTokens,
		BurnedTokens: ts.BurnedTokens,
		TotalSupply:  ts.TotalSupply,
	}
}

func registerCodec(cdc original.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "cosmos-sdk/MsgSend")
	cdc.RegisterConcrete(MsgMultiSend{}, "cosmos-sdk/MsgMultiSend")
}
