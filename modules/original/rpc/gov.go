package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
	"time"
)

type Gov interface {
	original.Module
	Deposit(proposalID uint64, amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error)
	Vote(proposalID uint64, option VoteOption, baseTx original.BaseTx) (original.ResultTx, original.Error)

	QueryProposal(proposalID uint64) (Proposal, original.Error)
	QueryProposals(request ProposalRequest) ([]Proposal, original.Error)

	QueryVote(proposalID uint64, voter string) (Vote, original.Error)
	QueryVotes(proposalID uint64) ([]Vote, original.Error)

	QueryDeposit(proposalID uint64, depositor string) (Deposit, original.Error)
	QueryDeposits(proposalID uint64) ([]Deposit, original.Error)

	QueryTally(proposalID uint64) (TallyResult, original.Error)
}

type VoteOption string

const (
	Yes        VoteOption = "Yes"
	No         VoteOption = "No"
	NoWithVeto VoteOption = "NoWithVeto"
	Abstain    VoteOption = "Abstain"
)

type Proposal interface {
	GetProposalID() string
	GetStatus() string
	GetTallyResult() TallyResult
	GetSubmitTime() time.Time
	GetDepositEndTime() time.Time
	GetTotalDeposit() original.Coins
	GetVotingStartTime() time.Time
	GetVotingEndTime() time.Time
}

var _ Proposal = (*BasicProposal)(nil)

type BasicProposal struct {
	ProposalID      string         `json:"proposal_id"` //  ID of the proposal
	Status          string         `json:"status"`
	TallyResult     TallyResult    `json:"tally_result"`      // Result of Tallys
	SubmitTime      time.Time      `json:"submit_time"`       // Time of the block where TxGovSubmitProposal was included
	DepositEndTime  time.Time      `json:"deposit_end_time"`  // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit    original.Coins `json:"total_deposit"`     // Current deposit on this proposal. Initial value is set at InitialDeposit
	VotingStartTime time.Time      `json:"voting_start_time"` // Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time      `json:"voting_end_time"`
}

func (b BasicProposal) GetProposalID() string {
	return b.ProposalID
}

func (b BasicProposal) GetStatus() string {
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

func (b BasicProposal) GetTotalDeposit() original.Coins {
	return b.TotalDeposit
}

func (b BasicProposal) GetVotingStartTime() time.Time {
	return b.VotingStartTime
}

func (b BasicProposal) GetVotingEndTime() time.Time {
	return b.VotingEndTime
}

// TallyResult defines a standard tally for a proposal
type TallyResult struct {
	Yes               string `json:"yes"`
	Abstain           string `json:"abstain"`
	No                string `json:"no"`
	NoWithVeto        string `json:"no_with_veto"`
	SystemVotingPower string `json:"system_voting_power,omitempty"`
}

type PlainTextProposal struct {
	Proposal
}

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
	Amount     original.Coins
}
