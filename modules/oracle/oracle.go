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
		Logger:         ac.Logger(),
	}
}

//CreateFeed create a stopped feed
func (o oracleClient) CreateFeed(request rpc.FeedCreateRequest) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(request.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		p, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	amt, err := o.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
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
		ServiceFeeCap:     amt,
		RepeatedFrequency: request.RepeatedFrequency,
		AggregateFunc:     request.AggregateFunc,
		ValueJsonPath:     request.ValueJsonPath,
		ResponseThreshold: request.ResponseThreshold,
	}
	return o.BuildAndSend([]sdk.Msg{msg}, request.BaseTx)
}

//StartFeed start a stopped feed
func (o oracleClient) StartFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgStartFeed{
		FeedName: feedName,
		Creator:  creator,
	}
	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//CreateAndStartFeed create and start a stopped feed
func (o oracleClient) CreateAndStartFeed(request rpc.FeedCreateRequest) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(request.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		p, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	amt, err := o.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
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
		ServiceFeeCap:     amt,
		RepeatedFrequency: request.RepeatedFrequency,
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
func (o oracleClient) PauseFeed(feedName string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(baseTx.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgPauseFeed{
		FeedName: feedName,
		Creator:  creator,
	}
	return o.BuildAndSend([]sdk.Msg{msg}, baseTx)
}

//EditFeed edit a feed
func (o oracleClient) EditFeed(request rpc.FeedEditRequest) (sdk.ResultTx, sdk.Error) {
	creator, err := o.QueryAddress(request.From)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	var providers []sdk.AccAddress
	for _, provider := range request.Providers {
		p, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return sdk.ResultTx{}, sdk.Wrapf("%s invalid address", p)
		}
		providers = append(providers, p)
	}

	amt, err := o.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return sdk.ResultTx{}, sdk.Wrap(err)
	}

	msg := MsgEditFeed{
		FeedName:          request.FeedName,
		LatestHistory:     request.LatestHistory,
		Description:       request.Description,
		Creator:           creator,
		Providers:         providers,
		Timeout:           request.Timeout,
		ServiceFeeCap:     amt,
		RepeatedFrequency: request.RepeatedFrequency,
		ResponseThreshold: request.ResponseThreshold,
	}
	return o.BuildAndSend([]sdk.Msg{msg}, request.BaseTx)
}

//QueryFeed return the feed by feedName
func (o oracleClient) QueryFeed(feedName string) (rpc.FeedContext, sdk.Error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}

	var ctx feedContext
	if err := o.QueryWithResponse("custom/oracle/feed", param, &ctx); err != nil {
		return rpc.FeedContext{}, sdk.Wrap(err)
	}
	return ctx.Convert().(rpc.FeedContext), nil
}

//QueryFeeds return all feeds by state
func (o oracleClient) QueryFeeds(state string) ([]rpc.FeedContext, sdk.Error) {
	param := struct {
		State string
	}{
		State: state,
	}

	var fcs feedContexts
	if err := o.QueryWithResponse("custom/oracle/feeds", param, &fcs); err != nil {
		return nil, sdk.Wrap(err)
	}
	return fcs.Convert().([]rpc.FeedContext), nil
}

//QueryFeedValue return all feed values by feedName
func (o oracleClient) QueryFeedValue(feedName string) ([]rpc.FeedValue, sdk.Error) {
	param := struct {
		FeedName string
	}{
		FeedName: feedName,
	}

	var fvs feedValues
	if err := o.QueryWithResponse("custom/oracle/feedValue", param, &fvs); err != nil {
		return nil, sdk.Wrap(err)
	}
	return fvs.Convert().([]rpc.FeedValue), nil
}

func (o oracleClient) RegisterFeedListener(feedName string, handler func(value rpc.FeedValue)) sdk.Error {
	feed, err := o.QueryFeed(feedName)
	if err != nil {
		return err
	}

	isInValidState := func(state string) bool {
		if state == COMPLETED || state == PAUSED || state == "" {
			return true
		}
		return false
	}

	if isInValidState(feed.State) {
		return sdk.Wrapf("feed:%s state is invalid:%s", feedName, feed.State)
	}

	handleResult := func(value string, sub1, sub2 sdk.Subscription) {
		o.Info().Str("feed-value", value).
			Msg("received feed value")
		var fv feedValue
		if err := cdc.UnmarshalJSON([]byte(value), &fv); err == nil {
			handler(fv.Convert().(rpc.FeedValue))
			f, err := o.QueryFeed(feedName)
			if err != nil || isInValidState(f.State) {
				_ = o.Unsubscribe(sub1)
				_ = o.Unsubscribe(sub2)
			}
		}
	}

	var sub1, sub2 sdk.Subscription

	blockBuilder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(tagFeedName).Contains(sdk.EventValue(feedName)))
	sub1, err = o.SubscribeNewBlock(blockBuilder, func(block sdk.EventDataNewBlock) {
		tagValue := tagFeedValue(feedName)
		result := block.ResultEndBlock.Tags.GetValue(tagValue)

		handleResult(result, sub1, sub2)
	})

	txBuilder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(tagFeedName).Contains(sdk.EventValue(feedName))).
		AddCondition(sdk.Cond(sdk.ActionKey).EQ("respond_service"))
	sub2, err = o.SubscribeTx(txBuilder, func(tx sdk.EventDataTx) {
		tagValue := tagFeedValue(feedName)
		result := tx.Result.Tags.GetValue(tagValue)

		handleResult(result, sub1, sub2)
	})
	return err
}
