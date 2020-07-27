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
	votes, err := gts.Gov().QueryVotes(2)
	require.NoError(gts.T(), err)
	fmt.Println(votes)
}

func (gts *GovTestSuite) TestQueryProposal() {
	proposal, err := gts.Gov().QueryProposal(1)
	fmt.Println(proposal)
	require.NoError(gts.T(), err)
	fmt.Println(proposal)
}
