package gov

import (
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type govClient struct {
	sdk.AbstractClient
	*log.Logger
}

func New(ac sdk.AbstractClient) sdk.Gov {
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
	return g.Broadcast(baseTx, []sdk.Msg{msg})
}

//Vote is responsible for voting for proposal
func (g govClient) Vote(proposalID uint64, option sdk.VoteOption, baseTx sdk.BaseTx) (sdk.Result, error) {
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
	return g.Broadcast(baseTx, []sdk.Msg{msg})
}

func (g govClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (g govClient) Name() string {
	return ModuleName
}
