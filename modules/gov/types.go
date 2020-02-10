package gov

import "github.com/irisnet/irishub-sdk-go/types"

type Gov interface {
	Deposit(proposalID uint64, amount types.Coins, baseTx types.BaseTx) (types.Result, error)
	Vote(proposalID uint64, option types.VoteOption, baseTx types.BaseTx) (types.Result, error)
	QueryDeposit(proposalID uint64, depositor string) (types.Deposit, error)
	QueryDeposits(proposalID uint64) (types.Deposits, error)
	QueryProposal(proposalID uint64) (types.Proposal, error)
	QueryProposals(params types.QueryProposalsParams) (types.Proposals, error)
	QueryVote(proposalID uint64, voter string) (types.Vote, error)
	QueryVotes(proposalID uint64) (types.Votes, error)
}

type govClient struct {
	types.AbstractClient
}

// Params for query 'custom/gov/deposit'
type QueryDepositParams struct {
	ProposalID uint64
	Depositor  types.AccAddress
}

// Params for query 'custom/gov/deposits'
type QueryDepositsParams struct {
	ProposalID uint64
}

// Params for query 'custom/gov/votes'
type QueryVotesParams struct {
	ProposalID uint64
}

// Params for query 'custom/gov/vote'
type QueryVoteParams struct {
	ProposalID uint64
	Voter      types.AccAddress
}

// Params for query 'custom/gov/proposal'
type QueryProposalParams struct {
	ProposalID uint64
}

// Params for query 'custom/gov/proposals'
type QueryProposalsParams struct {
	Voter          types.AccAddress
	Depositor      types.AccAddress
	ProposalStatus string
	Limit          uint64
}
