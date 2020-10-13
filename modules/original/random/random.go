package random

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"

	"github.com/irisnet/irishub-sdk-go/utils/log"
)

type randomClient struct {
	original.BaseClient
	*log.Logger
}

func Create(ac original.BaseClient) rpc.Random {
	return randomClient{
		BaseClient: ac,
		Logger:     ac.Logger(),
	}
}

func (r randomClient) RegisterCodec(cdc original.Codec) {
	registerCodec(cdc)
}

func (r randomClient) Name() string {
	return ModuleName
}

// Request is responsible for requesting a random number and callback `callback`
func (r randomClient) Request(request rpc.RandomRequest, baseTx original.BaseTx) (string, original.Error) {
	consumer, err := r.QueryAddress(baseTx.From)
	if err != nil {
		return "", original.Wrap(err)
	}

	amt, err := r.ToMinCoin(request.ServiceFeeCap...)
	if err != nil {
		return "", original.Wrap(err)
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
		baseTx.Mode = original.Commit
	}
	result, err := r.BuildAndSend([]original.Msg{msg}, baseTx)
	if err != nil {
		return "", original.Wrap(err)
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
	callback rpc.EventRequestRandomCallback) (original.Subscription, original.Error) {
	unsubscribe := func(sub1, sub2 original.Subscription) {
		_ = r.Unsubscribe(sub1)
		_ = r.Unsubscribe(sub2)
	}

	var sub1, sub2 original.Subscription

	blockBuilder := original.NewEventQueryBuilder().
		AddCondition(original.Cond(tagRequestID).Contains(original.EventValue(requestID)))
	sub1, err := r.SubscribeNewBlock(blockBuilder, func(block original.EventDataNewBlock) {
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

	txBuilder := original.NewEventQueryBuilder().
		AddCondition(original.Cond(tagRequestID).Contains(original.EventValue(requestID))).
		AddCondition(original.Cond(original.ActionKey).EQ("respond_service"))
	sub2, err = r.SubscribeTx(txBuilder, func(tx original.EventDataTx) {
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
func (r randomClient) QueryRandom(reqID string) (rpc.ResponseRandom, original.Error) {
	param := struct {
		ReqID string
	}{
		ReqID: reqID,
	}

	var rand rand
	if err := r.QueryWithResponse("custom/rand/rand", param, &rand); err != nil {
		return rpc.ResponseRandom{}, original.Wrap(err)
	}
	return rand.Convert().(rpc.ResponseRandom), nil
}

// QueryRequests returns the list of request by the specified block height
func (r randomClient) QueryRequests(height int64) ([]rpc.RequestRandom, original.Error) {
	param := struct {
		Height int64
	}{
		Height: height,
	}

	var rs requests
	if err := r.QueryWithResponse("custom/rand/queue", param, &rs); err != nil {
		return nil, original.Wrap(err)
	}
	return rs.Convert().([]rpc.RequestRandom), nil
}
