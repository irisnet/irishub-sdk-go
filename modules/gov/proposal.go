package gov

import (
	"time"

	"github.com/irisnet/irishub-sdk-go/types/rpc"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

var (
	_ Proposal = (*BasicProposal)(nil)
	_ Proposal = (*PlainTextProposal)(nil)
	_ Proposal = (*ParameterProposal)(nil)
	_ Proposal = (*CommunityTaxUsageProposal)(nil)
	_ Proposal = (*SoftwareUpgradeProposal)(nil)
)

// Proposal interface
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
	GetProposer() sdk.AccAddress
	ToSDKResponse() rpc.Proposal
}

// Basic Proposals
type BasicProposal struct {
	ProposalID      uint64         `json:"proposal_id"`       //  ID of the proposal
	Title           string         `json:"title"`             //  Title of the proposal
	Description     string         `json:"description"`       //  Description of the proposal
	ProposalType    string         `json:"proposal_type"`     //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Status          string         `json:"proposal_status"`   //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult     TallyResult    `json:"tally_result"`      //  Result of Tallys
	SubmitTime      time.Time      `json:"submit_time"`       //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime  time.Time      `json:"deposit_end_time"`  // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit    sdk.Coins      `json:"total_deposit"`     //  Current deposit on this proposal. Initial value is set at InitialDeposit
	VotingStartTime time.Time      `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time      `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
	Proposer        sdk.AccAddress `json:"proposer"`
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

func (b BasicProposal) GetTotalDeposit() sdk.Coins {
	return b.TotalDeposit
}

func (b BasicProposal) GetVotingStartTime() time.Time {
	return b.VotingStartTime
}

func (b BasicProposal) GetVotingEndTime() time.Time {
	return b.VotingEndTime
}

func (b BasicProposal) GetProposer() sdk.AccAddress {
	return b.Proposer
}

func (b BasicProposal) ToSDKResponse() rpc.Proposal {
	return rpc.BasicProposal{
		Title:          b.Title,
		Description:    b.Description,
		ProposalID:     b.ProposalID,
		ProposalStatus: b.Status,
		ProposalType:   b.ProposalType,
		TallyResult: rpc.TallyResult{
			Yes:               b.TallyResult.Yes,
			Abstain:           b.TallyResult.Abstain,
			No:                b.TallyResult.No,
			NoWithVeto:        b.TallyResult.NoWithVeto,
			SystemVotingPower: b.TallyResult.SystemVotingPower,
		},
		SubmitTime:      b.SubmitTime,
		DepositEndTime:  b.DepositEndTime,
		TotalDeposit:    b.TotalDeposit,
		VotingStartTime: b.VotingStartTime,
		VotingEndTime:   b.VotingEndTime,
		Proposer:        b.Proposer.String(),
	}
}

func (b BasicProposal) GetProposalID() uint64 {
	return b.ProposalID
}

type PlainTextProposal struct {
	BasicProposal
}

func (b PlainTextProposal) ToSDKResponse() rpc.Proposal {
	return rpc.PlainTextProposal{
		Proposal: b.BasicProposal.ToSDKResponse(),
	}
}

type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Params []Param

// Implements Proposal Interface
type ParameterProposal struct {
	BasicProposal
	Params Params `json:"params"`
}

func (b ParameterProposal) ToSDKResponse() rpc.Proposal {
	var params []rpc.Param
	for _, p := range b.Params {
		params = append(params, rpc.Param{
			Subspace: "", //TODO
			Key:      p.Key,
			SubKey:   "", //TODO
			Value:    p.Value,
		})
	}
	return rpc.ParameterProposal{
		Proposal: b.BasicProposal.ToSDKResponse(),
		Params:   params,
	}
}

// Implements Proposal Interface
type TaxUsage struct {
	Usage       string         `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Percent     string         `json:"percent"`
	Amount      sdk.Coins      `json:"amount"`
}

type CommunityTaxUsageProposal struct {
	BasicProposal
	TaxUsage TaxUsage `json:"tax_usage"`
}

func (b CommunityTaxUsageProposal) ToSDKResponse() rpc.Proposal {
	return rpc.CommunityTaxUsageProposal{
		Proposal: b.BasicProposal.ToSDKResponse(),
		TaxUsage: rpc.TaxUsage{
			Usage:       b.TaxUsage.Usage,
			DestAddress: b.TaxUsage.DestAddress.String(),
			Percent:     b.TaxUsage.Percent,
		},
	}
}

type SoftwareUpgradeProposal struct {
	BasicProposal
	ProtocolDefinition ProtocolDefinition `json:"protocol_definition"`
}

type ProtocolDefinition struct {
	Version   uint64 `json:"version"`
	Software  string `json:"software"`
	Height    uint64 `json:"height"`
	Threshold string `json:"threshold"`
}

func (b SoftwareUpgradeProposal) ToSDKResponse() rpc.Proposal {
	return rpc.SoftwareUpgradeProposal{
		Proposal: b.BasicProposal.ToSDKResponse(),
		ProtocolDefinition: rpc.ProtocolDefinition{
			Version:   b.ProtocolDefinition.Version,
			Software:  b.ProtocolDefinition.Software,
			Height:    b.ProtocolDefinition.Height,
			Threshold: b.ProtocolDefinition.Threshold,
		},
	}
}
