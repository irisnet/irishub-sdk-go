package oracle

import (
	"context"

	"github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
)

type oracleClient struct {
	sdk.BaseClient
	codec.Marshaler
}

func NewClient(baseClient sdk.BaseClient, marshaler codec.Marshaler) Client {
	return oracleClient{
		BaseClient: baseClient,
		Marshaler:  marshaler,
	}
}

func (oc oracleClient) Name() string {
	return ModuleName
}

func (oc oracleClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (oc oracleClient) CreateFeed(request CreateFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := oc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	serviceFeeCap, e := oc.ToMinCoin(request.ServiceFeeCap...)
	if e != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgCreateFeed{
		FeedName:          request.FeedName,
		LatestHistory:     request.LatestHistory,
		Description:       request.Description,
		Creator:           sender.String(),
		ServiceName:       request.ServiceName,
		Providers:         request.Providers,
		Input:             request.Input,
		Timeout:           request.Timeout,
		ServiceFeeCap:     serviceFeeCap,
		RepeatedFrequency: request.RepeatedFrequency,
		AggregateFunc:     request.AggregateFunc,
		ValueJsonPath:     request.ValueJsonPath,
		ResponseThreshold: request.ResponseThreshold,
	}
	return oc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (oc oracleClient) StartFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := oc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgStartFeed{
		FeedName: feedName,
		Creator:  sender.String(),
	}
	return oc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (oc oracleClient) PauseFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := oc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgPauseFeed{
		FeedName: feedName,
		Creator:  sender.String(),
	}
	return oc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (oc oracleClient) EditFeed(request EditFeedRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	sender, err := oc.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	serviceFeeCap, e := oc.ToMinCoin(request.ServiceFeeCap...)
	if e != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := &MsgEditFeed{
		FeedName:          request.FeedName,
		Description:       request.Description,
		LatestHistory:     request.LatestHistory,
		Providers:         request.Providers,
		Timeout:           request.Timeout,
		ServiceFeeCap:     serviceFeeCap,
		RepeatedFrequency: request.RepeatedFrequency,
		ResponseThreshold: request.ResponseThreshold,
		Creator:           sender.String(),
	}
	return oc.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

func (oc oracleClient) QueryFeed(feedName string) (QueryFeedResp, sdk.Error) {
	if len(feedName) == 0 {
		return QueryFeedResp{}, sdk.Wrapf("feedName is required")
	}

	conn, err := oc.GenConn()

	if err != nil {
		return QueryFeedResp{}, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Feed(
		context.Background(),
		&QueryFeedRequest{FeedName: feedName},
	)
	if err != nil {
		return QueryFeedResp{}, sdk.Wrap(err)
	}
	return res.Feed.Convert().(QueryFeedResp), nil
}

func (oc oracleClient) QueryFeeds(state string) ([]QueryFeedResp, sdk.Error) {
	// todo state (whether state is required)
	if len(state) == 0 {
		return nil, sdk.Wrapf("state is required")
	}

	conn, err := oc.GenConn()

	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).Feeds(
		context.Background(),
		&QueryFeedsRequest{State: state},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return feedContexts(res.Feeds).Convert().([]QueryFeedResp), nil
}

func (oc oracleClient) QueryFeedValue(feedName string) ([]QueryFeedValueResp, sdk.Error) {
	if len(feedName) == 0 {
		return nil, sdk.Wrapf("feedName is required")
	}

	conn, err := oc.GenConn()

	if err != nil {
		return nil, sdk.Wrap(err)
	}

	res, err := NewQueryClient(conn).FeedValue(
		context.Background(),
		&QueryFeedValueRequest{FeedName: feedName},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return feedValues(res.FeedValues).Convert().([]QueryFeedValueResp), nil
}
