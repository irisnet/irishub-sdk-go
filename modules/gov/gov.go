// Package gov provides governance functionalities
//
// [More Details](https://www.irisnet.org/docs/features/governance.html)
//
package gov

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type govClient struct {
	sdk.AbstractClient
	*log.Logger
}

func Create(ac sdk.AbstractClient) rpc.Gov {
	return govClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

//Deposit is responsible for depositing some tokens for proposal
func (g govClient) Deposit(proposalID uint64, amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error) {
	depositor, err := g.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	msg := MsgDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
	}
	g.Info().
		Uint64("proposalID", proposalID).
		Str("depositor", depositor.String()).
		Str("amount", amount.String()).
		Msg("execute gov deposit")
	return g.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//Vote is responsible for voting for proposal
func (g govClient) Vote(proposalID uint64, option rpc.VoteOption, baseTx sdk.BaseTx) (sdk.Result, error) {
	voter, err := g.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	op, err := VoteOptionFromString(option)
	if err != nil {
		return nil, err
	}

	msg := MsgVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     op,
	}
	g.Info().
		Uint64("proposalID", proposalID).
		Str("voter", voter.String()).
		Str("option", string(option)).
		Msg("execute gov vote")
	return g.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// QueryProposal returns the proposal of the specified proposalID
func (g govClient) QueryProposal(proposalID uint64) (rpc.Proposal, error) {
	param := struct {
		ProposalID uint64
	}{
		ProposalID: proposalID,
	}

	res, err := g.Query("custom/gov/proposal", param)
	if err != nil {
		return nil, err
	}

	var proposal proposal
	if err = cdc.UnmarshalJSON(res, &proposal); err != nil {
		return nil, err
	}

	return proposal.Convert().(rpc.Proposal), nil
}

// QueryProposals returns all proposals of the specified params
func (g govClient) QueryProposals(request rpc.ProposalRequest) (ps []rpc.Proposal, err error) {
	var voter, depositor sdk.AccAddress
	if len(request.Voter) != 0 {
		voter, err = sdk.AccAddressFromBech32(request.Voter)
		if err != nil {
			return nil, err
		}
	}

	if len(request.Depositor) != 0 {
		depositor, err = sdk.AccAddressFromBech32(request.Depositor)
		if err != nil {
			return nil, err
		}
	}

	param := struct {
		Voter          sdk.AccAddress
		Depositor      sdk.AccAddress
		ProposalStatus string
		Limit          uint64
	}{
		Voter:          voter,
		Depositor:      depositor,
		ProposalStatus: request.ProposalStatus,
		Limit:          request.Limit,
	}

	res, err := g.Query("custom/gov/proposals", param)
	if err != nil {
		return nil, err
	}

	var proposals proposals
	if err = cdc.UnmarshalJSON(res, &proposals); err != nil {
		return nil, err
	}
	for _, p := range proposals {
		ps = append(ps, p.Convert().(rpc.Proposal))
	}
	return ps, nil
}

// QueryVote returns the vote of the specified proposalID and voter
func (g govClient) QueryVote(proposalID uint64, voter string) (rpc.Vote, error) {
	v, err := sdk.AccAddressFromBech32(voter)
	if err != nil {
		return rpc.Vote{}, err
	}

	param := struct {
		ProposalID uint64
		Voter      sdk.AccAddress
	}{
		ProposalID: proposalID,
		Voter:      v,
	}

	var vote vote
	err = g.QueryWithResponse("custom/gov/vote", param, &vote)
	if err != nil {
		return rpc.Vote{}, err
	}
	return vote.Convert().(rpc.Vote), nil
}

// QueryVotes returns all votes of the specified proposalID
func (g govClient) QueryVotes(proposalID uint64) ([]rpc.Vote, error) {
	param := struct {
		ProposalID uint64
	}{
		ProposalID: proposalID,
	}

	var vs votes
	err := g.QueryWithResponse("custom/gov/votes", param, &vs)
	if err != nil {
		return nil, err
	}
	return vs.Convert().([]rpc.Vote), nil
}

// QueryDeposit returns the deposit of the specified proposalID and depositor
func (g govClient) QueryDeposit(proposalID uint64, depositor string) (rpc.Deposit, error) {
	d, err := sdk.AccAddressFromBech32(depositor)
	if err != nil {
		return rpc.Deposit{}, err
	}

	param := struct {
		ProposalID uint64
		Depositor  sdk.AccAddress
	}{
		ProposalID: proposalID,
		Depositor:  d,
	}

	var deposit deposit
	err = g.QueryWithResponse("custom/gov/deposit", param, &deposit)
	if err != nil {
		return rpc.Deposit{}, err
	}
	return deposit.Convert().(rpc.Deposit), nil
}

// QueryDeposits returns all deposits of the specified proposalID
func (g govClient) QueryDeposits(proposalID uint64) ([]rpc.Deposit, error) {
	param := struct {
		ProposalID uint64
	}{
		ProposalID: proposalID,
	}

	var deposits deposits
	err := g.QueryWithResponse("custom/gov/deposits", param, &deposits)
	if err != nil {
		return nil, err
	}
	return deposits.Convert().([]rpc.Deposit), nil
}

// QueryTally returns the result of proposal by the specified proposalID
func (g govClient) QueryTally(proposalID uint64) (rpc.TallyResult, error) {
	param := struct {
		ProposalID uint64
	}{
		ProposalID: proposalID,
	}

	var tally tallyResult
	err := g.QueryWithResponse("custom/gov/tally", param, &tally)
	if err != nil {
		return rpc.TallyResult{}, err
	}
	return tally.Convert().(rpc.TallyResult), nil
}

func (g govClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (g govClient) Name() string {
	return ModuleName
}
