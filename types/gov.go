package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/utils"
)

//-----------------------------------------------------------
// MsgDeposit
type MsgDeposit struct {
	ProposalID uint64     `json:"proposal_id"` // ID of the proposal
	Depositor  AccAddress `json:"depositor"`   // Address of the depositor
	Amount     Coins      `json:"amount"`      // Coins to add to the proposal's deposit
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Type() string { return "deposit" }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() error {
	if len(msg.Depositor) == 0 {
		return errors.New(fmt.Sprintf("account %s is invalid", msg.Depositor))
	}
	if msg.ProposalID < 0 {
		return errors.New(fmt.Sprintf("Unknown proposal with id %d", msg.ProposalID))
	}
	return nil
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	b, err := defaultCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []AccAddress {
	return []AccAddress{msg.Depositor}
}

//-----------------------------------------------------------
// MsgVote
type MsgVote struct {
	ProposalID uint64     `json:"proposal_id"` // ID of the proposal
	Voter      AccAddress `json:"voter"`       //  address of the voter
	Option     VoteOption `json:"option"`      //  option from OptionSet chosen by the voter
}

// Implements Msg.
// nolint
func (msg MsgVote) Type() string { return "vote" }

// Implements Msg.
func (msg MsgVote) ValidateBasic() error {
	if len(msg.Voter.Bytes()) == 0 {
		return errors.New(fmt.Sprintf("account %s is invalid", msg.Voter))
	}
	if msg.ProposalID < 0 {
		return errors.New(fmt.Sprintf("Unknown proposal with id %d", msg.ProposalID))
	}
	if !ValidVoteOption(msg.Option) {
		return errors.New(fmt.Sprintf("'%v' is not a valid voting option", msg.Option))
	}
	return nil
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	b, err := defaultCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return utils.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []AccAddress {
	return []AccAddress{msg.Voter}
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
func VoteOptionFromString(str string) (VoteOption, error) {
	switch str {
	case "Yes":
		return OptionYes, nil
	case "Abstain":
		return OptionAbstain, nil
	case "No":
		return OptionNo, nil
	case "NoWithVeto":
		return OptionNoWithVeto, nil
	default:
		return VoteOption(0xff), errors.New(fmt.Sprintf("'%s' is not a valid vote option", str))
	}
}

// Is defined VoteOption
func ValidVoteOption(option VoteOption) bool {
	if option == OptionYes ||
		option == OptionAbstain ||
		option == OptionNo ||
		option == OptionNoWithVeto {
		return true
	}
	return false
}

// Marshal needed for protobuf compatibility
func (vo VoteOption) Marshal() ([]byte, error) {
	return []byte{byte(vo)}, nil
}

// Unmarshal needed for protobuf compatibility
func (vo *VoteOption) Unmarshal(data []byte) error {
	*vo = VoteOption(data[0])
	return nil
}

// Marshals to JSON using string
func (vo VoteOption) MarshalJSON() ([]byte, error) {
	return json.Marshal(vo.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (vo *VoteOption) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := VoteOptionFromString(s)
	if err != nil {
		return err
	}
	*vo = bz2
	return nil
}

// Turns VoteOption byte to String
func (vo VoteOption) String() string {
	switch vo {
	case OptionYes:
		return "Yes"
	case OptionAbstain:
		return "Abstain"
	case OptionNo:
		return "No"
	case OptionNoWithVeto:
		return "NoWithVeto"
	default:
		return ""
	}
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (vo VoteOption) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(fmt.Sprintf("%s", vo.String())))
	default:
		s.Write([]byte(fmt.Sprintf("%v", byte(vo))))
	}
}

type VoteResult struct {
	Voter      string `json:"voter"`
	ProposalID string `json:"proposal_id"`
	Option     string `json:"option"`
}

type ProposalResult struct {
}
