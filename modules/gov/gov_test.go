package gov_test

import (
	"testing"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type GovTestSuite struct {
	suite.Suite
	sim.TestClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(GovTestSuite))
}

func (gts *GovTestSuite) SetupTest() {
	gts.TestClient = sim.NewClient()
}

func (gts *GovTestSuite) TestDeposit() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Memo: "test",
		Mode: sdk.Commit,
	}

	amt, _ := sdk.NewIntFromString("10000000000000000000000")
	amount := sdk.NewCoins(sdk.NewCoin("iris-atto", amt))
	proposalID := uint64(12)

	proposal, err := gts.Gov().QueryProposal(proposalID)
	require.True(gts.T(), err.IsNil())
	require.Equal(gts.T(), proposalID, proposal.GetProposalID())

	proposals, err := gts.Gov().QueryProposals(rpc.ProposalRequest{
		Depositor: gts.Sender().String(),
	})
	require.True(gts.T(), err.IsNil())
	require.NotEmpty(gts.T(), proposals)

	rs, err := gts.Gov().Deposit(proposalID, amount, baseTx)
	require.True(gts.T(), err.IsNil())
	require.True(gts.T(), rs.IsSuccess())

	d, err := gts.Gov().QueryDeposit(proposalID, gts.Sender().String())
	require.True(gts.T(), err.IsNil())
	require.NotEmpty(gts.T(), d)

	ds, err := gts.Gov().QueryDeposits(proposalID)
	require.True(gts.T(), err.IsNil())
	require.NotEmpty(gts.T(), ds)

	rs, err = gts.Gov().Vote(proposalID, rpc.Yes, baseTx)
	require.True(gts.T(), err.IsNil())
	require.True(gts.T(), rs.IsSuccess())

	vote, err := gts.Gov().QueryVote(proposalID, gts.Sender().String())
	require.True(gts.T(), err.IsNil())
	require.Equal(gts.T(), proposalID, vote.ProposalID)

	votes, err := gts.Gov().QueryVotes(proposalID)
	require.True(gts.T(), err.IsNil())
	require.NotEmpty(gts.T(), votes)

	tally, err := gts.Gov().QueryTally(proposalID)
	require.True(gts.T(), err.IsNil())
	require.NotEmpty(gts.T(), tally.Yes)
}
