package random

import (
	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type randomClient struct {
	sdk.AbstractClient
	*log.Logger
}

func Create(ac sdk.AbstractClient) rpc.Random {
	return randomClient{
		AbstractClient: ac,
		Logger:         ac.Logger(),
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

	amt, err := r.ConvertToMinCoin(request.ServiceFeeCap...)
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

	requestID := result.Tags.GetValue(tagRequestID)
	if needWatch && !request.Oracle {
		var subscription sdk.Subscription
		subscription, err = r.SubscribeNewBlockWithQuery(nil, func(block sdk.EventDataNewBlock) {
			tags := block.ResultBeginBlock.Tags
			r.Debug().
				Int64("height", block.Block.Height).
				Str("tags", tags.String()).
				Msg("received block")
			requestIDs := tags.GetValues(tagRequestID)
			for _, reqID := range requestIDs {
				if reqID == requestID {
					result, err := r.QueryRandom(requestID)
					var randomNum string
					if err == nil {
						randomNum = result.Value
					}
					request.Callback(requestID, randomNum, err)
					_ = r.Unscribe(subscription)
					return
				}
			}
		})
	}
	return requestID, nil
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
