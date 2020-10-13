package asset

import (
	json2 "encoding/json"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"strconv"

	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/utils/json"
)

const (
	ModuleName = "token"
)

var (
	cdc = original.NewAminoCodec()

	_ original.Msg = &MsgIssueToken{}
	_ original.Msg = &MsgEditToken{}
	_ original.Msg = &MsgMintToken{}
	_ original.Msg = &MsgTransferTokenOwner{}
)

func init() {
	registerCodec(cdc)
}

// MsgIssueToken
type MsgIssueToken struct {
	Symbol        string              `json:"symbol"`
	Name          string              `json:"name"`
	Decimal       uint8               `json:"decimal"`
	MinUnitAlias  string              `json:"min_unit_alias"`
	InitialSupply uint64              `json:"initial_supply"`
	MaxSupply     uint64              `json:"max_supply"`
	Mintable      bool                `json:"mintable"`
	Owner         original.AccAddress `json:"owner"`
}

func (msg MsgIssueToken) Route() string { return ModuleName }

// Implements Msg.
func (msg MsgIssueToken) Type() string { return "issue_token" }

// Implements Msg.
func (msg MsgIssueToken) ValidateBasic() error {
	//nothing
	return nil
}

// Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// Implements Msg.
func (msg MsgIssueToken) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Owner}
}

// MsgTransferTokenOwner for transferring the token owner
type MsgTransferTokenOwner struct {
	SrcOwner original.AccAddress `json:"src_owner"` // the current owner address of the token
	DstOwner original.AccAddress `json:"dst_owner"` // the new owner
	Symbol   string              `json:"symbol"`    // the token symbol
}

// GetSignBytes implements Msg
func (msg MsgTransferTokenOwner) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// GetSigners implements Msg
func (msg MsgTransferTokenOwner) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.SrcOwner}
}

func (msg MsgTransferTokenOwner) ValidateBasic() error {
	//nothing
	return nil
}

func (msg MsgTransferTokenOwner) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgTransferTokenOwner) Type() string { return "transfer_token_owner" }

// MsgEditToken for editing a specified token
type MsgEditToken struct {
	Symbol    string              `json:"symbol"` //  symbol of token
	Owner     original.AccAddress `json:"owner"`  //  owner of token
	MaxSupply uint64              `json:"max_supply"`
	Mintable  Bool                `json:"mintable"` //  mintable of token
	Name      string              `json:"name"`
}

func (msg MsgEditToken) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgEditToken) Type() string { return "edit_token" }

// ValidateBasic implements Msg
func (msg MsgEditToken) ValidateBasic() error {
	//nothing
	return nil
}

// GetSignBytes implements Msg
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return json.MustSort(b)
}

// GetSigners implements Msg
func (msg MsgEditToken) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Owner}
}

// MsgMintToken for minting the token to a specified address
type MsgMintToken struct {
	Symbol string              `json:"symbol"` // the symbol of the token
	Owner  original.AccAddress `json:"owner"`  // the current owner address of the token
	To     original.AccAddress `json:"to"`     // address of minting token to
	Amount uint64              `json:"amount"` // amount of minting token
}

func (msg MsgMintToken) Route() string { return ModuleName }

// Type implements Msg
func (msg MsgMintToken) Type() string { return "mint_token" }

// GetSignBytes implements Msg
func (msg MsgMintToken) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// GetSigners implements Msg
func (msg MsgMintToken) GetSigners() []original.AccAddress {
	return []original.AccAddress{msg.Owner}
}

// ValidateBasic implements Msg
func (msg MsgMintToken) ValidateBasic() error {
	return nil
}

// tokenFees is for the token fees query output
type tokenFees struct {
	Exist    bool          `json:"exist"`     // indicate if the token has existed
	IssueFee original.Coin `json:"issue_fee"` // issue fee
	MintFee  original.Coin `json:"mint_fee"`  // mint fee
}

func (t tokenFees) Convert() interface{} {
	return rpc.TokenFees{
		Exist:    t.Exist,
		IssueFee: t.IssueFee,
		MintFee:  t.MintFee,
	}
}

type Bool string

func (b Bool) ToBool() bool {
	v := string(b)
	if len(v) == 0 {
		return false
	}
	result, _ := strconv.ParseBool(v)
	return result
}

func (b Bool) String() string {
	return string(b)
}

// Marshal needed for protobuf compatibility
func (b Bool) Marshal() ([]byte, error) {
	return []byte(b), nil
}

// Unmarshal needed for protobuf compatibility
func (b *Bool) Unmarshal(data []byte) error {
	*b = Bool(data[:])
	return nil
}

// Marshals to JSON using string
func (b Bool) MarshalJSON() ([]byte, error) {
	return json2.Marshal(b.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (b *Bool) UnmarshalJSON(data []byte) error {
	var s string
	err := json2.Unmarshal(data, &s)
	if err != nil {
		return nil
	}
	*b = Bool(s)
	return nil
}

// Token defines a struct for the fungible token
type Token struct {
	Symbol        string              `json:"symbol"`
	Name          string              `json:"name"`
	Scale         uint8               `json:"scale"`
	MinUnit       string              `json:"min_unit"`
	InitialSupply uint64              `json:"initial_supply"`
	MaxSupply     uint64              `json:"max_supply"`
	Mintable      bool                `json:"mintable"`
	Owner         original.AccAddress `json:"owner"`
}

func registerCodec(cdc original.Codec) {
	cdc.RegisterConcrete(MsgIssueToken{}, "irismod/token/MsgIssueToken")
	cdc.RegisterConcrete(MsgEditToken{}, "irismod/token/MsgEditToken")
	cdc.RegisterConcrete(MsgMintToken{}, "irismod/token/MsgMintToken")
	cdc.RegisterConcrete(MsgTransferTokenOwner{}, "irismod/token/MsgTransferTokenOwner")
	cdc.RegisterConcrete(Token{}, "irismod/token/Token")
}
