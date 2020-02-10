package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tendermint/go-amino"
	"time"

	"github.com/irisnet/irishub-sdk-go/utils"
)

var (
	_ Msg      = MsgDeposit{}
	_ Msg      = MsgVote{}
	_ Proposal = BasicProposal{}
	_ Proposal = (*PlainTextProposal)(nil)
	_ Proposal = (*ParameterProposal)(nil)
	_ Proposal = (*CommunityTaxUsageProposal)(nil)
	_ Proposal = (*SoftwareUpgradeProposal)(nil)
	_ Proposal = (*SystemHaltProposal)(nil)
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
	cdc := amino.NewCodec()
	b, err := cdc.MarshalJSON(msg)
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
	cdc := amino.NewCodec()
	b, err := cdc.MarshalJSON(msg)
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

// Deposit
type Deposit struct {
	Depositor  AccAddress `json:"depositor"`   //  Address of the depositor
	ProposalID uint64     `json:"proposal_id"` //  proposalID of the proposal
	Amount     Coins      `json:"amount"`      //  Deposit amount
}
type Deposits []Deposit

type Vote struct {
	Voter      AccAddress `json:"voter"`       //  address of the voter
	ProposalID uint64     `json:"proposal_id"` //  proposalID of the proposal
	Option     VoteOption `json:"option"`      //  option from OptionSet chosen by the voter
}
type Votes []Vote

// Tally Results
type TallyResult struct {
	Yes               Dec `json:"yes"`
	Abstain           Dec `json:"abstain"`
	No                Dec `json:"no"`
	NoWithVeto        Dec `json:"no_with_veto"`
	SystemVotingPower Dec `json:"system_voting_power"`
}
type Proposal interface {
	GetProposalID() uint64
	GetTitle() string
	GetDescription() string
	GetProposalType() string
	GetProposalStatus() string
	GetTallyResult() TallyResult
	GetSubmitTime() time.Time
	GetDepositEndTime() time.Time
	GetTotalDeposit() Coins
	GetVotingStartTime() time.Time
	GetVotingEndTime() time.Time
	GetProposer() AccAddress
}
type Proposals []Proposal

type BasicProposal struct {
	ProposalID      uint64      `json:"proposal_id"`       //  ID of the proposal
	Title           string      `json:"title"`             //  Title of the proposal
	Description     string      `json:"description"`       //  Description of the proposal
	ProposalType    string      `json:"proposal_type"`     //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Status          string      `json:"proposal_status"`   //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult     TallyResult `json:"tally_result"`      //  Result of Tallys
	SubmitTime      time.Time   `json:"submit_time"`       //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime  time.Time   `json:"deposit_end_time"`  // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit    Coins       `json:"total_deposit"`     //  Current deposit on this proposal. Initial value is set at InitialDeposit
	VotingStartTime time.Time   `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time   `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
	Proposer        AccAddress  `json:"proposer"`
}

func (b BasicProposal) GetProposalID() uint64 {
	return b.ProposalID
}

func (b BasicProposal) GetTitle() string {
	return b.Title
}

func (b BasicProposal) GetDescription() string {
	return b.Description
}

func (b BasicProposal) GetProposalType() string {
	return b.ProposalType
}

func (b BasicProposal) GetProposalStatus() string {
	return b.Status
}

func (b BasicProposal) GetTallyResult() TallyResult {
	return b.TallyResult
}

func (b BasicProposal) GetSubmitTime() time.Time {
	return b.SubmitTime
}

func (b BasicProposal) GetDepositEndTime() time.Time {
	return b.DepositEndTime
}

func (b BasicProposal) GetTotalDeposit() Coins {
	return b.TotalDeposit
}

func (b BasicProposal) GetVotingStartTime() time.Time {
	return b.VotingStartTime
}

func (b BasicProposal) GetVotingEndTime() time.Time {
	return b.VotingEndTime
}

func (b BasicProposal) GetProposer() AccAddress {
	return b.Proposer
}

// Implements Proposal Interface
type PlainTextProposal struct {
	BasicProposal
}

// Implements Proposal Interface
type SystemHaltProposal struct {
	BasicProposal
}

// Implements Proposal Interface
type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}
type Params []Param
type ParameterProposal struct {
	BasicProposal
	Params Params `json:"params"`
}

// Implements Proposal Interface
type TaxUsage struct {
	Usage       string     `json:"usage"`
	DestAddress AccAddress `json:"dest_address"`
	Percent     Dec        `json:"percent"`
	Amount      Coins      `json:"amount"`
}
type CommunityTaxUsageProposal struct {
	BasicProposal
	TaxUsage TaxUsage `json:"tax_usage"`
}

//Implements Proposal Interface
type ProtocolDefinition struct {
	Version   uint64 `json:"version"`
	Software  string `json:"software"`
	Height    uint64 `json:"height"`
	Threshold Dec    `json:"threshold"`
}
type SoftwareUpgradeProposal struct {
	BasicProposal
	ProtocolDefinition ProtocolDefinition `json:"protocol_definition"`
}

// Params for query 'custom/gov/proposals'
type QueryProposalsParams struct {
	Voter          string
	Depositor      string
	ProposalStatus string
	Limit          uint64
}

func RegisterGov(cdc *amino.Codec) {
	cdc.RegisterInterface((*Proposal)(nil), nil)
	cdc.RegisterConcrete(&BasicProposal{}, "irishub/gov/BasicProposal", nil)
	cdc.RegisterConcrete(&ParameterProposal{}, "irishub/gov/ParameterProposal", nil)
	cdc.RegisterConcrete(&PlainTextProposal{}, "irishub/gov/PlainTextProposal", nil)
	cdc.RegisterConcrete(&SoftwareUpgradeProposal{}, "irishub/gov/SoftwareUpgradeProposal", nil)
	cdc.RegisterConcrete(&SystemHaltProposal{}, "irishub/gov/SystemHaltProposal", nil)
	cdc.RegisterConcrete(&CommunityTaxUsageProposal{}, "irishub/gov/CommunityTaxUsageProposal", nil)
	cdc.RegisterConcrete(&Vote{}, "irishub/gov/Vote", nil)
	cdc.RegisterConcrete(MsgDeposit{}, "irishub/gov/MsgDeposit", nil)
	cdc.RegisterConcrete(MsgVote{}, "irishub/gov/MsgVote", nil)
}
