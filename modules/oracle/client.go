package oracle

import (
	"github.com/irisnet/irishub-sdk-go/types"
)

type oracleClient struct {
	types.AbstractClient
}

func New(ac types.AbstractClient) types.Oracle {
	return oracleClient{
		AbstractClient: ac,
	}
}

func (o oracleClient) CreateFeed(request types.FeedCreateRequest) (result types.Result, err error) {
	panic("implement me")
}

func (o oracleClient) StartFeed(feedName string) (result types.Result, err error) {
	panic("implement me")
}

func (o oracleClient) PauseFeed(feedName string) (result types.Result, err error) {
	panic("implement me")
}

func (o oracleClient) EditFeed(request types.FeedEditRequest) (result types.Result, err error) {
	panic("implement me")
}

func (o oracleClient) QueryFeed(feedName string) (feed types.FeedContext, err error) {
	panic("implement me")
}

func (o oracleClient) QueryFeeds(state string) (feed types.FeedContext, err error) {
	panic("implement me")
}

func (o oracleClient) QueryFeedValue(feedName string) (value []types.FeedValue, err error) {
	panic("implement me")
}
