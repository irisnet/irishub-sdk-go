package gov

import (
	"time"

	sdk "github.com/irisnet/core-sdk-go/types"
)

// expose Gov module api for user
type Client interface {
	sdk.Module
	SubmitProposal(request SubmitProposalRequest, baseTx sdk.BaseTx) (uint64, sdk.ResultTx, sdk.Error)
	Deposit(request DepositRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	Vote(request VoteRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryProposal(proposalId uint64) (QueryProposalResp, sdk.Error)
	QueryProposals(proposalStatus string) ([]QueryProposalResp, sdk.Error)
	QueryVote(proposalId uint64, voter string) (QueryVoteResp, sdk.Error)
	QueryVotes(proposalId uint64) ([]QueryVoteResp, sdk.Error)
	QueryParams(paramsType string) (QueryParamsResp, sdk.Error)
	QueryDeposit(proposalId uint64, depositor string) (QueryDepositResp, sdk.Error)
	QueryDeposits(proposalId uint64) ([]QueryDepositResp, sdk.Error)
	QueryTallyResult(proposalId uint64) (QueryTallyResultResp, sdk.Error)
}

type SubmitProposalRequest struct {
	Title          string       `json:"title"`
	Description    string       `json:"description"`
	Type           string       `json:"type"`
	InitialDeposit sdk.DecCoins `json:"initial_deposit"`
}

type DepositRequest struct {
	ProposalId uint64       `json:"proposal_id"`
	Amount     sdk.DecCoins `json:"amount"`
}

type VoteRequest struct {
	ProposalId uint64 `json:"proposal_id"`
	Option     string `json:"option"`
}

type QueryProposalResp struct {
	ProposalId       uint64               `json:"proposal_id"`
	Content          Content              `json:"content"`
	Status           string               `json:"status"`
	FinalTallyResult QueryTallyResultResp `json:"final_tally_result"`
	SubmitTime       time.Time            `json:"submit_time"`
	DepositEndTime   time.Time            `json:"deposit_end_time"`
	TotalDeposit     sdk.Coins            `json:"total_deposit"`
	VotingStartTime  time.Time            `json:"voting_start_time"`
	VotingEndTime    time.Time            `json:"voting_end_time"`
}

type QueryVoteResp struct {
	ProposalId uint64 `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     int32  `json:"option"`
}

type (
	votingParams struct {
		VotingPeriod time.Duration `json:"voting_period"`
	}
	depositParams struct {
		MinDeposit       sdk.Coins     `json:"min_deposit"`
		MaxDepositPeriod time.Duration `json:"max_deposit_period"`
	}
	tallyParams struct {
		Quorum        sdk.Dec `json:"quorum"`
		Threshold     sdk.Dec `json:"threshold"`
		VetoThreshold sdk.Dec `json:"veto_threshold"`
	}
	QueryParamsResp struct {
		VotingParams  votingParams  `json:"voting_params"`
		DepositParams depositParams `json:"deposit_params"`
		TallyParams   tallyParams   `json:"tally_params"`
	}
)

type QueryDepositResp struct {
	ProposalId uint64    `json:"proposal_id"`
	Depositor  string    `json:"depositor"`
	Amount     sdk.Coins `json:"amount"`
}

type QueryTallyResultResp struct {
	Yes        sdk.Int `json:"yes"`
	Abstain    sdk.Int `json:"abstain"`
	No         sdk.Int `json:"no"`
	NoWithVeto sdk.Int `json:"no_with_veto"`
}
