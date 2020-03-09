package oracle

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type oracleClient struct {
	sdk.AbstractClient
	*log.Logger
}

func (o oracleClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (o oracleClient) Name() string {
	return ModuleName
}

func Create(ac sdk.AbstractClient) rpc.Oracle {
	return oracleClient{
		AbstractClient: ac,
		Logger:         ac.Logger().With(ModuleName),
	}
}

//CreateFeed create a stopped feed
func (o oracleClient) CreateFeed(request rpc.FeedCreateRequest) (result sdk.Result, err error) {
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
	return o.BuildAndSend([]sdk.Msg{msg}, request.BaseTx)
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
	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//CreateAndStartFeed create and start a stopped feed
func (o oracleClient) CreateAndStartFeed(request rpc.FeedCreateRequest) (result sdk.Result, err error) {
	creator, err := o.QueryAddress(request.From, request.Password)
	if err != nil {
		return nil, err
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		providers = append(providers, sdk.MustAccAddressFromBech32(provider))
	}

	msgCreateFeed := MsgCreateFeed{
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

	msgStartFeed := MsgStartFeed{
		FeedName: request.FeedName,
		Creator:  creator,
	}
	return o.BuildAndSend([]sdk.Msg{msgCreateFeed, msgStartFeed}, request.BaseTx)
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
	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//EditFeed edit a feed
func (o oracleClient) EditFeed(request rpc.FeedEditRequest) (result sdk.Result, err error) {
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
	return o.BuildAndSend([]sdk.Msg{msg}, request.BaseTx)
}

//QueryFeed return the feed by feedName
func (o oracleClient) QueryFeed(feedName string) (rpc.FeedContext, error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}

	var ctx feedContext
	if err := o.QueryWithResponse("custom/oracle/feed", param, &ctx); err != nil {
		return rpc.FeedContext{}, err
	}
	return ctx.Convert().(rpc.FeedContext), nil
}

//QueryFeeds return all feeds by state
func (o oracleClient) QueryFeeds(state string) ([]rpc.FeedContext, error) {
	param := struct {
		State string
	}{
		State: state,
	}

	var fcs feedContexts
	if err := o.QueryWithResponse("custom/oracle/feeds", param, &fcs); err != nil {
		return nil, err
	}
	return fcs.Convert().([]rpc.FeedContext), nil
}

//QueryFeedValue return all feed values by feedName
func (o oracleClient) QueryFeedValue(feedName string) ([]rpc.FeedValue, error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}

	var fvs feedValues
	if err := o.QueryWithResponse("custom/oracle/feedValue", param, &fvs); err != nil {
		return nil, err
	}
	return fvs.Convert().([]rpc.FeedValue), nil
}
