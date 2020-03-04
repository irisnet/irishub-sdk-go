package gov_test

import (
	"testing"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	amt, _ := sdk.NewIntFromString("10000000000000000000000")
	amount := sdk.NewCoins(sdk.NewCoin("iris-atto", amt))
	proposalID := uint64(6)

	proposal, err := gts.Gov().QueryProposal(proposalID)
	require.NoError(gts.T(), err)
	require.Equal(gts.T(), proposalID, proposal.GetProposalID())

	proposals, err := gts.Gov().QueryProposals(sdk.ProposalRequest{
		Depositor: gts.Sender().String(),
	})
	require.NoError(gts.T(), err)
	require.NotEmpty(gts.T(), proposals)

	rs, err := gts.Gov().Deposit(proposalID, amount, baseTx)
	require.NoError(gts.T(), err)
	require.True(gts.T(), rs.IsSuccess())

	d, err := gts.Gov().QueryDeposit(proposalID, gts.Sender().String())
	require.NoError(gts.T(), err)
	require.NotEmpty(gts.T(), d)

	ds, err := gts.Gov().QueryDeposits(proposalID)
	require.NoError(gts.T(), err)
	require.NotEmpty(gts.T(), ds)

	rs, err = gts.Gov().Vote(proposalID, sdk.Yes, baseTx)
	require.NoError(gts.T(), err)
	require.True(gts.T(), rs.IsSuccess())

	vote, err := gts.Gov().QueryVote(proposalID, gts.Sender().String())
	require.NoError(gts.T(), err)
	require.Equal(gts.T(), proposalID, vote.ProposalID)

	votes, err := gts.Gov().QueryVotes(proposalID)
	require.NoError(gts.T(), err)
	require.NotEmpty(gts.T(), votes)

	tally, err := gts.Gov().QueryTally(proposalID)
	require.NoError(gts.T(), err)
	require.NotEmpty(gts.T(), tally.Yes)
}
