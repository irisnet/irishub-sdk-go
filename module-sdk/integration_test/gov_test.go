package integrationtest

import (
	"encoding/json"
	"fmt"

	"github.com/irisnet/core-sdk-go/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/gov-sdk-go"
)

func (s IntegrationTestSuite) TestGov() {
	cases := []SubTest{
		{
			"TestGov",
			testGov,
		},

		{
			"TestParams",
			testParams,
		},
	}

	for _, t := range cases {
		s.Run(t.testName, func() {
			t.testCase(s)
		})
	}

}

func testGov(s IntegrationTestSuite) {
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "TEST",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	// send submitProposal tx
	submitProposalReq := gov.SubmitProposalRequest{
		Title:       s.RandStringOfLength(4),
		Description: s.RandStringOfLength(6),
		Type:        "Text",
	}
	proposalId, res, err := s.Gov.SubmitProposal(submitProposalReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	// query proposal details based on ProposalID.
	proposal, err := s.Gov.QueryProposal(proposalId)
	require.NoError(s.T(), err)
	require.Equal(s.T(), proposal.ProposalId, proposalId)
	require.Equal(s.T(), "PROPOSAL_STATUS_DEPOSIT_PERIOD", proposal.Status)

	// query all proposals based on given status.
	proposalStatus := proposal.Status
	proposals, err := s.Gov.QueryProposals(proposalStatus)
	var exists bool
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), proposals)
	for _, proposal := range proposals {
		if proposal.ProposalId == proposalId {
			exists = true
		}
	}
	require.True(s.T(), exists)

	// send Deposit tx
	amount, e := types.ParseDecCoins("2000iris")
	require.NoError(s.T(), e)
	depositReq := gov.DepositRequest{
		ProposalId: proposalId,
		Amount:     amount,
	}
	res, err = s.Gov.Deposit(depositReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	// query single deposit information based proposalID, depositAddr.
	depositor := s.Account().Address.String()
	deposit, err := s.Gov.QueryDeposit(proposalId, depositor)
	require.NoError(s.T(), err)
	require.Equal(s.T(), "2000000000uiris", deposit.Amount.String())

	// query all deposits of a single proposal.
	deposits, err := s.Gov.QueryDeposits(proposalId)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), deposits)
	for _, deposit := range deposits {
		require.Equal(s.T(), proposalId, deposit.ProposalId)
	}

	// send vote tx
	voteReq := gov.VoteRequest{
		ProposalId: proposalId,
		Option:     "VOTE_OPTION_YES",
	}
	res, err = s.Gov.Vote(voteReq, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	// query voted information based on proposalID, voterAddr.
	voter := s.Account().Address.String()
	vote, err := s.Gov.QueryVote(proposalId, voter)
	require.NoError(s.T(), err)
	require.Equal(s.T(), proposalId, vote.ProposalId)

	// query votes of a given proposal.
	votes, err := s.Gov.QueryVotes(proposalId)
	require.NoError(s.T(), err)
	require.Greater(s.T(), len(votes), 0)

	// query the tally of a proposal vote.
	_, err = s.Gov.QueryTallyResult(proposalId)
	require.NoError(s.T(), err)
}

func testParams(s IntegrationTestSuite) {
	paramsTypes := []string{"voting", "tallying", "deposit"}
	for _, paramType := range paramsTypes {
		res, err := s.Gov.QueryParams(paramType)
		require.NoError(s.T(), err)
		bz, e := json.Marshal(res)
		require.NoError(s.T(), e)
		fmt.Println(string(bz))
	}
}
