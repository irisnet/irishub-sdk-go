package oracle

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type oracleClient struct {
	sdk.AbstractClient
}

func New(ac sdk.AbstractClient) sdk.Oracle {
	return oracleClient{
		AbstractClient: ac,
	}
}

//CreateFeed create a stopped feed
func (o oracleClient) CreateFeed(request sdk.FeedCreateRequest) (result sdk.Result, err error) {
	creator, err := o.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		providers = append(providers, sdk.MustAccAddressFromBech32(provider))
	}

	msg := MsgCreateFeed{
		FeedName:          request.FeedName,
		LatestHistory:     request.LatestHistory,
		Description:       request.Description,
		Creator:           creator,
		ServiceName:       request.ServiceName,
		Providers:         providers,
		Input:             request.Input,
		Timeout:           request.Timeout,
		ServiceFeeCap:     request.ServiceFeeCap,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		AggregateFunc:     request.AggregateFunc,
		ValueJsonPath:     request.ValueJsonPath,
		ResponseThreshold: request.ResponseThreshold,
	}
	return o.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

//StartFeed start a stopped feed
func (o oracleClient) StartFeed(feedName string, baseTx sdk.BaseTx) (result sdk.Result, err error) {
	creator, err := o.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	msg := MsgStartFeed{
		FeedName: feedName,
		Creator:  creator,
	}
	return o.Broadcast(baseTx, []sdk.Msg{msg})
}

//PauseFeed pause a running feed
func (o oracleClient) PauseFeed(feedName string, baseTx sdk.BaseTx) (result sdk.Result, err error) {
	creator, err := o.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, err
	}

	msg := MsgPauseFeed{
		FeedName: feedName,
		Creator:  creator,
	}
	return o.Broadcast(baseTx, []sdk.Msg{msg})
}

//EditFeed edit a feed
func (o oracleClient) EditFeed(request sdk.FeedEditRequest) (result sdk.Result, err error) {
	creator, err := o.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		providers = append(providers, sdk.MustAccAddressFromBech32(provider))
	}

	msg := MsgEditFeed{
		FeedName:          request.FeedName,
		LatestHistory:     request.LatestHistory,
		Description:       request.Description,
		Creator:           creator,
		Providers:         providers,
		Timeout:           request.Timeout,
		ServiceFeeCap:     request.ServiceFeeCap,
		RepeatedFrequency: request.RepeatedFrequency,
		RepeatedTotal:     request.RepeatedTotal,
		ResponseThreshold: request.ResponseThreshold,
	}
	return o.Broadcast(request.BaseTx, []sdk.Msg{msg})
}

//QueryFeed return the feed by feedName
func (o oracleClient) QueryFeed(feedName string) (feed sdk.FeedContext, err error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}
	var ctx FeedContext
	err = o.Query("custom/oracle/feed", param, &ctx)
	if err != nil {
		return feed, err
	}
	return ctx.toSDKFeedContext(), nil
}

//QueryFeeds return all feeds by state
func (o oracleClient) QueryFeeds(state string) (feed []sdk.FeedContext, err error) {
	param := struct {
		State string
	}{
		State: state,
	}
	var ctx FeedContexts
	err = o.Query("custom/oracle/feeds", param, &ctx)
	if err != nil {
		return feed, err
	}
	return ctx.toSDKFeedContexts(), nil
}

//QueryFeedValue return all feed values by feedName
func (o oracleClient) QueryFeedValue(feedName string) (value []sdk.FeedValue, err error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}
	err = o.Query("custom/oracle/feedValue", param, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}
