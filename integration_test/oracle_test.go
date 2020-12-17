package integration_test

import (
	"github.com/irisnet/irishub-sdk-go/modules/oracle"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func (s IntegrationTestSuite) TestOracle() {
	// todo complete oracle module
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	input := ``

	feedName := s.RandStringOfLength(6)
	createFeedRequest := oracle.CreateFeedRequest{
		FeedName:          feedName,
		LatestHistory:     0,
		Description:       s.RandStringOfLength(10),
		ServiceName:       s.RandStringOfLength(6),
		Providers:         nil,
		Input:             input,
		Timeout:           0,
		ServiceFeeCap:     nil,
		RepeatedFrequency: 0,
		AggregateFunc:     "",
		ValueJsonPath:     "",
		ResponseThreshold: 0,
	}
	res, err := s.Oracle.CreateFeed(createFeedRequest, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	pauseFeedRequest := oracle.PauseFeedRequest{
		FeedName: feedName,
		Creator:  "",
	}
	res, err = s.Oracle.PauseFeed(pauseFeedRequest, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	startFeedRequest := oracle.StartFeedRequest{
		FeedName: feedName,
		Creator:  "",
	}
	res, err = s.Oracle.StartFeed(startFeedRequest, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)

	editFeedRequest := oracle.EditFeedRequest{
		FeedName:          feedName,
		Description:       s.RandStringOfLength(10),
		LatestHistory:     0,
		Providers:         nil,
		Timeout:           0,
		ServiceFeeCap:     nil,
		RepeatedFrequency: 0,
		ResponseThreshold: 0,
	}
	res, err = s.Oracle.EditFeed(editFeedRequest, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.Hash)
}
