package rpc

import (
	"time"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Gov interface {
	sdk.Module
	Deposit(proposalID uint64, amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error)
	Vote(proposalID uint64, option VoteOption, baseTx sdk.BaseTx) (sdk.Result, error)

	QueryProposal(proposalID uint64) (Proposal, error)
	QueryProposals(request ProposalRequest) ([]Proposal, error)

	QueryVote(proposalID uint64, voter string) (Vote, error)
	QueryVotes(proposalID uint64) ([]Vote, error)

	QueryDeposit(proposalID uint64, depositor string) (Deposit, error)
	QueryDeposits(proposalID uint64) ([]Deposit, error)

	QueryTally(proposalID uint64) (TallyResult, error)
}

type VoteOption string

const (
	Yes        VoteOption = "Yes"
	No         VoteOption = "No"
	NoWithVeto VoteOption = "NoWithVeto"
	Abstain    VoteOption = "Abstain"
)

//=========================basicProposal========================================================
type Proposal interface {
	GetProposalID() uint64
	GetTitle() string
	GetDescription() string
	GetProposalType() string
	GetStatus() string
	GetTallyResult() TallyResult
	GetSubmitTime() time.Time
	GetDepositEndTime() time.Time
	GetTotalDeposit() sdk.Coins
	GetVotingStartTime() time.Time
	GetVotingEndTime() time.Time
	GetProposer() string
}

var _ Proposal = (*BasicProposal)(nil)

type BasicProposal struct {
	ProposalID      uint64      `json:"proposal_id"` //  ID of the proposal
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	ProposalType    string      `json:"proposal_type"`
	ProposalStatus  string      `json:"proposal_status"`   // Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult     TallyResult `json:"tally_result"`      // Result of Tallys
	SubmitTime      time.Time   `json:"submit_time"`       // Time of the block where TxGovSubmitProposal was included
	DepositEndTime  time.Time   `json:"deposit_end_time"`  // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit    sdk.Coins   `json:"total_deposit"`     // Current deposit on this proposal. Initial value is set at InitialDeposit
	VotingStartTime time.Time   `json:"voting_start_time"` // Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time   `json:"voting_end_time"`
	Proposer        string      `json:"proposer"`
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

func (b BasicProposal) GetStatus() string {
	return b.ProposalStatus
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

func (b BasicProposal) GetTotalDeposit() sdk.Coins {
	return b.TotalDeposit
}

func (b BasicProposal) GetVotingStartTime() time.Time {
	return b.VotingStartTime
}

func (b BasicProposal) GetVotingEndTime() time.Time {
	return b.VotingEndTime
}

func (b BasicProposal) GetProposer() string {
	return b.Proposer
}

// TallyResult defines a standard tally for a proposal
type TallyResult struct {
	Yes               string `json:"yes"`
	Abstain           string `json:"abstain"`
	No                string `json:"no"`
	NoWithVeto        string `json:"no_with_veto"`
	SystemVotingPower string `json:"system_voting_power,omitempty"`
}

//==========================PlainTextProposal=======================================================
type PlainTextProposal struct {
	Proposal
}

//==========================ParameterProposal=======================================================
type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	SubKey   string `json:"sub_key,omitempty"`
	Value    string `json:"value"`
}

type ParameterProposal struct {
	Proposal
	Params []Param `json:"params"`
}

//==========================CommunityTaxUsageProposal===============================================
type TaxUsage struct {
	Usage       string `json:"usage"`
	DestAddress string `json:"dest_address"`
	Percent     string `json:"percent"`
	//Amount      Coins  `json:"amount"`
}
type CommunityTaxUsageProposal struct {
	Proposal
	TaxUsage TaxUsage `json:"tax_usage"`
}

//==========================SoftwareUpgradeProposal===============================================
type ProtocolDefinition struct {
	Version   uint64 `json:"version"`
	Software  string `json:"software"`
	Height    uint64 `json:"height"`
	Threshold string `json:"threshold"`
}

type SoftwareUpgradeProposal struct {
	Proposal
	ProtocolDefinition
}

//============================for query=================================================================
type ProposalRequest struct {
	Voter          string
	Depositor      string
	ProposalStatus string
	Limit          uint64
}

type Vote struct {
	Voter      string
	ProposalID uint64
	Option     string
}

type Deposit struct {
	Depositor  string
	ProposalID uint64
	Amount     sdk.Coins
}
