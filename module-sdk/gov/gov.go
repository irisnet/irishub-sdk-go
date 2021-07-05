package gov

import (
	"context"
	"strconv"

	"github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/query"
)

type govClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(baseClient sdk.BaseClient, marshaler codec.Marshaler) Client {
	return govClient{
		BaseClient: baseClient,
		Marshaler:  marshaler,
	}
}

func (gc govClient) Name() string {
	return ModuleName
}

func (gc govClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (gc govClient) SubmitProposal(request SubmitProposalRequest, baseTx sdk.BaseTx) (uint64, sdk.ResultTx, sdk.Error) {
	proposer, err := gc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return 0, sdk.ResultTx{}, sdk.Wrap(err)
	}

	deposit, err := gc.ToMinCoin(request.InitialDeposit...)
	if err != nil {
		return 0, sdk.ResultTx{}, sdk.Wrap(err)
	}

	content := ContentFromProposalType(request.Title, request.Description, request.Type)
	msg, e := NewMsgSubmitProposal(content, deposit, proposer)
	if e != nil {
		return 0, sdk.ResultTx{}, sdk.Wrap(err)
	}

	result, err := gc.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return 0, sdk.ResultTx{}, sdk.Wrap(err)
	}

	proposalIdStr, e := result.Events.GetValue(sdk.EventTypeSubmitProposal, AttributeKeyProposalId)
	if e != nil {
		return 0, result, sdk.Wrap(e)
	}

	proposalId, e := strconv.Atoi(proposalIdStr)
	if e != nil {
		return 0, result, sdk.Wrap(e)
	}
	return uint64(proposalId), result, err
}

func (gc govClient) Deposit(request DepositRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	depositor, err := gc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	amount, err := gc.ToMinCoin(request.Amount...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgDeposit{
		ProposalId: request.ProposalId,
		Depositor:  depositor.String(),
		Amount:     amount,
	}
	return gc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

// about VoteRequest.Option see  VoteOption_value
func (gc govClient) Vote(request VoteRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	voter, err := gc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	option := VoteOption_value[request.Option]
	msg := &MsgVote{
		ProposalId: request.ProposalId,
		Voter:      voter.String(),
		Option:     VoteOption(option),
	}
	return gc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (gc govClient) QueryProposal(proposalId uint64) (QueryProposalResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryProposalResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Proposal(
		context.Background(),
		&QueryProposalRequest{
			ProposalId: proposalId,
		})
	if err != nil {
		return QueryProposalResp{}, sdk.Wrap(err)
	}
	return res.Proposal.Convert().(QueryProposalResp), nil
}

// if proposalStatus is nil will return all status's proposals
// about proposalStatus see VoteOption_value
func (gc govClient) QueryProposals(proposalStatus string) ([]QueryProposalResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Proposals(
		context.Background(),
		&QueryProposalsRequest{
			ProposalStatus: ProposalStatus(VoteOption_value[proposalStatus]),
			Pagination: &query.PageRequest{
				Offset:     0,
				Limit:      100,
				CountTotal: true,
			},
		})
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return Proposals(res.Proposals).Convert().([]QueryProposalResp), nil
}

// about QueryVoteResp.Option see VoteOption_name
func (gc govClient) QueryVote(proposalId uint64, voter string) (QueryVoteResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryVoteResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Vote(
		context.Background(),
		&QueryVoteRequest{
			ProposalId: proposalId,
			Voter:      voter,
		})
	if err != nil {
		return QueryVoteResp{}, sdk.Wrap(err)
	}
	return res.Vote.Convert().(QueryVoteResp), nil
}

func (gc govClient) QueryVotes(proposalId uint64) ([]QueryVoteResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Votes(
		context.Background(),
		&QueryVotesRequest{
			ProposalId: proposalId,
			Pagination: &query.PageRequest{
				Offset:     0,
				Limit:      100,
				CountTotal: true,
			},
		})
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return Votes(res.Votes).Convert().([]QueryVoteResp), nil
}

// QueryParams params_type("voting", "tallying", "deposit"), if don't pass will return all params_typ res
func (gc govClient) QueryParams(paramsType string) (QueryParamsResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Params(
		context.Background(),
		&QueryParamsRequest{
			ParamsType: paramsType,
		},
	)
	if err != nil {
		return QueryParamsResp{}, sdk.Wrap(err)
	}
	return res.Convert().(QueryParamsResp), nil
}

func (gc govClient) QueryDeposit(proposalId uint64, depositor string) (QueryDepositResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryDepositResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Deposit(
		context.Background(),
		&QueryDepositRequest{
			ProposalId: proposalId,
			Depositor:  depositor,
		},
	)
	if err != nil {
		return QueryDepositResp{}, sdk.Wrap(err)
	}
	return res.Deposit.Convert().(QueryDepositResp), nil
}

func (gc govClient) QueryDeposits(proposalId uint64) ([]QueryDepositResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Deposits(
		context.Background(),
		&QueryDepositsRequest{
			ProposalId: proposalId,
			Pagination: &query.PageRequest{
				Offset:     0,
				Limit:      100,
				CountTotal: true,
			},
		},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return Deposits(res.Deposits).Convert().([]QueryDepositResp), nil
}

func (gc govClient) QueryTallyResult(proposalId uint64) (QueryTallyResultResp, sdk.Error) {
	conn, err := gc.GenConn()

	if err != nil {
		return QueryTallyResultResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).TallyResult(
		context.Background(),
		&QueryTallyResultRequest{
			ProposalId: proposalId,
		},
	)
	if err != nil {
		return QueryTallyResultResp{}, sdk.Wrap(err)
	}
	return res.Tally.Convert().(QueryTallyResultResp), nil
}
