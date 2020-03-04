package gov

import (
	"errors"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/tools/json"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "gov"
)

var (
	_ sdk.Msg = MsgDeposit{}
	_ sdk.Msg = MsgVote{}

	cdc = sdk.NewAminoCodec()
)

//-----------------------------------------------------------
// MsgDeposit
type MsgDeposit struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Depositor  sdk.AccAddress `json:"depositor"`   // Address of the depositor
	Amount     sdk.Coins      `json:"amount"`      // Coins to add to the proposal's deposit
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Type() string { return "deposit" }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() error {
	if len(msg.Depositor) == 0 {
		return errors.New("depositor is empty")
	}
	if msg.ProposalID < 0 {
		return errors.New("invalid proposalID")
	}
	return nil
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

//-----------------------------------------------------------
// MsgVote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	Option     VoteOption     `json:"option"`      //  option from OptionSet chosen by the voter
}

// Implements Msg.
// nolint
func (msg MsgVote) Type() string { return "vote" }

// Implements Msg.
func (msg MsgVote) ValidateBasic() error {
	if len(msg.Voter) == 0 {
		return errors.New("voter is empty")
	}
	if msg.ProposalID < 0 {
		return errors.New("invalid proposalID")
	}
	if msg.Option != OptionYes &&
		msg.Option != OptionNo &&
		msg.Option != OptionNoWithVeto &&
		msg.Option != OptionAbstain {
		return errors.New("invalid option")
	}
	return nil
}

func (msg MsgVote) String() string {
	return fmt.Sprintf("MsgVote{%v - %s}", msg.ProposalID, msg.Option)
}

// Implements Msg.
func (msg MsgVote) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return json.MustSort(b)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}

// Type that represents VoteOption as a byte
type VoteOption byte

//nolint
const (
	OptionEmpty      VoteOption = 0x00
	OptionYes        VoteOption = 0x01
	OptionAbstain    VoteOption = 0x02
	OptionNo         VoteOption = 0x03
	OptionNoWithVeto VoteOption = 0x04
)

// String to proposalType byte.  Returns ff if invalid.
func VoteOptionFromString(option sdk.VoteOption) (VoteOption, error) {
	switch option {
	case sdk.Yes:
		return OptionYes, nil
	case sdk.Abstain:
		return OptionAbstain, nil
	case sdk.No:
		return OptionNo, nil
	case sdk.NoWithVeto:
		return OptionNoWithVeto, nil
	default:
		return OptionEmpty, errors.New(fmt.Sprintf("'%s' is not a valid vote option", str))
	}
}

func registerCodec(cdc sdk.Codec) {
	cdc.RegisterConcrete(MsgDeposit{}, "irishub/gov/MsgDeposit")
	cdc.RegisterConcrete(MsgVote{}, "irishub/gov/MsgVote")
}
