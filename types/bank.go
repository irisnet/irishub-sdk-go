package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/utils"
	"github.com/tendermint/go-amino"
)

type Bank interface {
	GetAccount(address string) (BaseAccount, error)
	GetTokenStats(tokenID string) (TokenStats, error)
	Send(to string, amount Coins, baseTx BaseTx) (Result, error)
	Burn(amount Coins, baseTx BaseTx) (Result, error)
	SetMemoRegexp(memoRegexp string, baseTx BaseTx) (Result, error)
}

type MsgSend struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

var _ Msg = MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(in []Input, out []Output) MsgSend {
	return MsgSend{Inputs: in, Outputs: out}
}

// Implements Msg.
func (msg MsgSend) Type() string  { return "send" }

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

// MsgSetMemoRegexp - set memo regexp
type MsgSetMemoRegexp struct {
	Owner      AccAddress `json:"owner"`
	MemoRegexp string     `json:"memo_regexp"`
}

type TokenStats struct {
	LooseTokens  Coins `json:"loose_tokens"`
	BondedTokens Coins `json:"bonded_tokens"`
	BurnedTokens Coins `json:"burned_tokens"`
	TotalSupply  Coins `json:"total_supply"`
}

// defines the params for query: "custom/acc/account"
type QueryAccountParams struct {
	Address AccAddress
}

// QueryTokenParams is the query parameters for 'custom/asset/tokens/{id}'
type QueryTokenParams struct {
	TokenId string
}

func RegisterBank(cdc *amino.Codec) {
	cdc.RegisterConcrete(MsgSend{}, "irishub/bank/Send", nil)
}
