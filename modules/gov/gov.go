package gov

import (
	"github.com/irisnet/irishub-sdk-go/types"
)

func New(ac types.AbstractClient) Gov {
	return govClient{
		AbstractClient: ac,
	}
}

func (g govClient) Deposit(proposalID uint64, amount types.Coins, baseTx types.BaseTx) (types.Result, error) {
	depositor := g.GetSender(baseTx.From)
	msg := types.MsgDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
	}
	return g.Broadcast(baseTx, []types.Msg{msg})
}

func (g govClient) Vote(proposalID uint64, option types.VoteOption, baseTx types.BaseTx) (types.Result, error) {
	voter := g.GetSender(baseTx.From)
	msg := types.MsgVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
	return g.Broadcast(baseTx, []types.Msg{msg})
}

func (g govClient) QueryDeposit(proposalID uint64, depositor string) (result types.Deposit, err error) {
	addr, err := types.AccAddressFromBech32(depositor)
	if err != nil {
		return result, err
	}
	param := QueryDepositParams{
		ProposalID: proposalID,
		Depositor:  addr,
	}

	err = g.Query("custom/gov/deposit", param, &result)
	return result, err
}

func (g govClient) QueryDeposits(proposalID uint64) (result types.Deposits, err error) {
	param := QueryDepositsParams{
		ProposalID: proposalID,
	}

	err = g.Query("custom/gov/deposits", param, &result)
	return result, err
}

func (g govClient) QueryProposal(proposalID uint64) (result types.Proposal, err error) {
	param := QueryProposalParams{
		ProposalID: proposalID,
	}

	err = g.Query("custom/gov/proposal", param, &result)
	return result, err
}

func (g govClient) QueryProposals(params types.QueryProposalsParams) (result types.Proposals, err error) {
	var p QueryProposalsParams
	if len(params.Depositor) > 0 {
		depositor, err := types.AccAddressFromBech32(params.Depositor)
		if err != nil {
			return result, err
		}
		p.Depositor = depositor
	}

	if len(params.Voter) > 0 {
		voter, err := types.AccAddressFromBech32(params.Voter)
		if err != nil {
			return result, err
		}
		p.Voter = voter
	}
	p.Limit = params.Limit
	p.ProposalStatus = params.ProposalStatus
	err = g.Query("custom/gov/proposals", p, &result)
	return result, err
}

func (g govClient) QueryVote(proposalID uint64, voter string) (result types.Vote, err error) {
	addr, err := types.AccAddressFromBech32(voter)
	if err != nil {
		return result, err
	}

	param := QueryVoteParams{
		ProposalID: proposalID,
		Voter:      addr,
	}

	err = g.Query("custom/gov/vote", param, &result)
	return result, err
}

func (g govClient) QueryVotes(proposalID uint64) (result types.Votes, err error) {
	param := QueryVotesParams{
		ProposalID: proposalID,
	}

	err = g.Query("custom/gov/votes", param, &result)
	return result, err
}
