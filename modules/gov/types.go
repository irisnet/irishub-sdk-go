package gov

import "github.com/irisnet/irishub-sdk-go/types"

type Gov interface {
	Deposit(proposalID uint64, amount types.Coins, baseTx types.BaseTx) (types.Result, error)
	Vote(proposalID uint64, option types.VoteOption, baseTx types.BaseTx) (types.Result, error)
	QueryDeposit(proposalID uint64, depositor types.AccAddress) (types.VoteResult, error)
	QueryDeposits(proposalID uint64) (types.VoteResult, error)
	QueryProposal(proposalID uint64) (types.ProposalResult, error)
	QueryProposals(proposalID uint64) (types.ProposalResult, error)
	QueryVote(proposalID uint64, voter types.AccAddress) (types.VoteResult, error)
	QueryVotes(proposalID uint64) (types.VoteResult, error)
}

type govClient struct {
	types.AbstractClient
}

func (g govClient) Deposit(proposalID uint64, amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}

func (g govClient) Vote(proposalID uint64, option types.VoteOption, baseTx types.BaseTx) (types.Result, error) {
	panic("implement me")
}

func (g govClient) QueryDeposit(proposalID uint64, depositor types.AccAddress) (types.VoteResult, error) {
	panic("implement me")
}

func (g govClient) QueryDeposits(proposalID uint64) (types.VoteResult, error) {
	panic("implement me")
}

func (g govClient) QueryProposal(proposalID uint64) (types.ProposalResult, error) {
	panic("implement me")
}

func (g govClient) QueryProposals(proposalID uint64) (types.ProposalResult, error) {
	panic("implement me")
}

func (g govClient) QueryVote(proposalID uint64, voter types.AccAddress) (types.VoteResult, error) {
	panic("implement me")
}

func (g govClient) QueryVotes(proposalID uint64) (types.VoteResult, error) {
	panic("implement me")
}
