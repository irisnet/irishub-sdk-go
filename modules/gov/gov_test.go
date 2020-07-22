package gov_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/suite"
)

type GovTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(GovTestSuite))
}

func (gts *GovTestSuite) SetupTest() {
	gts.MockClient = test.GetMock()
}

func (gts *GovTestSuite) TestQueryVoters() {
	votes, err := gts.Gov().QueryVotes(123)
	require.NoError(gts.T(), err)
	fmt.Println(votes)
}

//func (gts *GovTestSuite) TestDeposit() {
//	baseTx := sdk.BaseTx{
//		From: "test1",
//		gas:  20000,
//		Memo: "test",
//		mode: sdk.Commit,
//	}
//
//	amt, _ := sdk.NewIntFromString("10000000000000000000000")
//	amount := sdk.NewCoins(sdk.NewCoin("iris-atto", amt))
//	proposalID := uint64(12)
//
//	proposal, err := gts.Gov().QueryProposal(proposalID)
//	require.NoError(gts.T(), err)
//	require.EQ(gts.T(), proposalID, proposal.GetProposalID())
//
//	proposals, err := gts.Gov().QueryProposals(rpc.ProposalRequest{
//		Depositor: gts.Sender().String(),
//	})
//	require.NoError(gts.T(), err)
//	require.NotEmpty(gts.T(), proposals)
//
//	rs, err := gts.Gov().Deposit(proposalID, amount, baseTx)
//	require.NoError(gts.T(), err)
//	require.NotEmpty(gts.T(), rs.Hash)
//
//	d, err := gts.Gov().QueryDeposit(proposalID, gts.Sender().String())
//	require.NoError(gts.T(), err)
//	require.NotEmpty(gts.T(), d)
//
//	ds, err := gts.Gov().QueryDeposits(proposalID)
//	require.NoError(gts.T(), err)
//	require.NotEmpty(gts.T(), ds)
//
//	rs, err = gts.Gov().Vote(proposalID, rpc.Yes, baseTx)
//	require.NoError(gts.T(), err)
//	require.NotEmpty(gts.T(), rs.Hash)
//
//	vote, err := gts.Gov().QueryVote(proposalID, gts.Sender().String())
//	require.NoError(gts.T(), err)
//	require.EQ(gts.T(), proposalID, vote.ProposalID)
//
//	tally, err := gts.Gov().QueryTally(proposalID)
//	require.NoError(gts.T(), err)
//	require.NotEmpty(gts.T(), tally.Yes)
//}
