package gov_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/gov"
	"github.com/irisnet/irishub-sdk-go/types"
	"testing"

	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GovTestSuite struct {
	types.Codec
	suite.Suite
	gov.Gov
	sender, validator string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(GovTestSuite))
}

func (gts *GovTestSuite) SetupTest() {
	tc := sim.NewTestClient()
	gts.Gov = tc.Gov
	gts.sender = tc.GetTestSender()
}

func (gts GovTestSuite) TestDeposit() {
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	amt := types.NewIntWithDecimal(100, 18)
	coin := types.NewCoin("iris-atto", amt)
	res, err := gts.Deposit(1, types.NewCoins(coin), baseTx)
	require.NoError(gts.T(), err)
	require.True(gts.T(), res.IsSuccess())
}

func (gts GovTestSuite) TestVote() {
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	res, err := gts.Vote(1, types.OptionYes, baseTx)
	require.NoError(gts.T(), err)
	require.True(gts.T(), res.IsSuccess())
}

func (gts GovTestSuite) TestQueryDeposit() {
	deposit, err := gts.QueryDeposit(1, gts.sender)
	require.NoError(gts.T(), err)
	fmt.Printf("%v", deposit)
}

func (gts GovTestSuite) TestQueryDeposits() {
	deposit, err := gts.QueryDeposits(1)
	require.NoError(gts.T(), err)
	fmt.Printf("%v", deposit)
}

func (gts GovTestSuite) TestQueryProposal() {
	p, err := gts.QueryProposal(1)
	require.NoError(gts.T(), err)
	bz, _ := types.NewAmino().MarshalJSON(p)
	fmt.Printf("%s", string(bz))
}

func (gts GovTestSuite) TestQueryProposals() {
	p, err := gts.QueryProposals(types.QueryProposalsParams{})
	require.NoError(gts.T(), err)
	bz, _ := types.NewAmino().MarshalJSON(p)
	fmt.Printf("%s", string(bz))
}

func (gts GovTestSuite) TestQueryVote() {
	p, err := gts.QueryVote(1, gts.sender)
	require.NoError(gts.T(), err)
	bz, _ := types.NewAmino().MarshalJSON(p)
	fmt.Printf("%s", string(bz))
}

func (gts GovTestSuite) TestQueryVotes() {
	p, err := gts.QueryVotes(1)
	require.NoError(gts.T(), err)
	bz, _ := types.NewAmino().MarshalJSON(p)
	fmt.Printf("%s", string(bz))
}
