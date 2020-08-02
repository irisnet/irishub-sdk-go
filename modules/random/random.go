package random

import (
	"github.com/irisnet/irishub-sdk-go/rpc"

	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type randomClient struct {
	sdk.BaseClient
	*log.Logger
}

func Create(ac sdk.BaseClient) rpc.Random {
	return randomClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (r randomClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (r randomClient) Name() string {
	return ModuleName
}

// Request is responsible for requesting a random number and callback `callback`
func (r randomClient) Request(request rpc.RandomRequest, baseTx sdk.BaseTx) (string, sdk.Error) {
	consumer, err := r.QueryAddress(baseTx.From)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	amt, err := r.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	msg := MsgRequestRand{
		Consumer:      consumer,
		BlockInterval: request.BlockInterval,
		Oracle:        request.Oracle,
		ServiceFeeCap: amt,
	}

	needWatch := request.Callback != nil
	if needWatch {
		//mode must be set to commit
		baseTx.Mode = sdk.Commit
	}
	result, err := r.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	requestID := result.Events.GetValues(tagRequestID, "")
	if needWatch {
		_, err := r.SubscribeRandom(requestID[0], request.Callback)
		r.Err(err).
			Str(tagRequestID, requestID[0]).
			Msg("subscribe random failed")
	}
	return requestID[0], nil
}

func (r randomClient) SubscribeRandom(requestID string,
	callback rpc.EventRequestRandomCallback) (sdk.Subscription, sdk.Error) {
	unsubscribe := func(sub1, sub2 sdk.Subscription) {
		_ = r.Unsubscribe(sub1)
		_ = r.Unsubscribe(sub2)
	}

	var sub1, sub2 sdk.Subscription

	blockBuilder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(tagRequestID).Contains(sdk.EventValue(requestID)))
	sub1, err := r.SubscribeNewBlock(blockBuilder, func(block sdk.EventDataNewBlock) {
		//cancel subscribe
		unsubscribe(sub1, sub2)

		rand := block.ResultBeginBlock.Events.GetValues(tagRand(tagRequestID), "")
		r.Debug().
			Int64("height", block.Block.Height).
			Str("requestID", requestID).
			Str("random", rand[0]).
			Msg("received random result")

		callback(requestID, rand[0], nil)
	})

	txBuilder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(tagRequestID).Contains(sdk.EventValue(requestID))).
		AddCondition(sdk.Cond(sdk.ActionKey).EQ("respond_service"))
	sub2, err = r.SubscribeTx(txBuilder, func(tx sdk.EventDataTx) {
		//cancel subscribe
		unsubscribe(sub1, sub2)

		rand := tx.Result.Events.GetValues(tagRand(tagRequestID), "")
		r.Debug().
			Int64("height", tx.Height).
			Str("requestID", requestID).
			Str("random", rand[0]).
			Msg("received random result")

		callback(requestID, rand[0], nil)
	})
	return sub1, err
}

// QueryRandom returns the random information of the specified reqID
func (r randomClient) QueryRandom(reqID string) (rpc.ResponseRandom, sdk.Error) {
	param := struct {
		ReqID string
	}{
		ReqID: reqID,
	}

	var rand rand
	if err := r.QueryWithResponse("custom/rand/rand", param, &rand); err != nil {
		return rpc.ResponseRandom{}, sdk.Wrap(err)
	}
	return rand.Convert().(rpc.ResponseRandom), nil
}

// QueryRequests returns the list of request by the specified block height
func (r randomClient) QueryRequests(height int64) ([]rpc.RequestRandom, sdk.Error) {
	param := struct {
		Height int64
	}{
		Height: height,
	}

	var rs requests
	if err := r.QueryWithResponse("custom/rand/queue", param, &rs); err != nil {
		return nil, sdk.Wrap(err)
	}
	return rs.Convert().([]rpc.RequestRandom), nil
}
