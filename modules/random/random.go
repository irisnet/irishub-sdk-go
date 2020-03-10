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
		Logger:         ac.Logger().With(ModuleName),
	}
}

func (r randomClient) RegisterCodec(cdc sdk.Codec) {
	registerCodec(cdc)
}

func (r randomClient) Name() string {
	return ModuleName
}

// Generate is responsible for requesting a random number and callback `callback`
func (r randomClient) Generate(request rpc.RandomRequest) (string, sdk.Error) {
	consumer, err := r.QueryAddress(request.From, request.Password)
	if err != nil {
		return "", sdk.Wrap(err)
	}

	msg := MsgRequestRand{
		Consumer:      consumer,
		BlockInterval: request.BlockInterval,
	}

	//mode must be set to commit
	request.BaseTx.Mode = sdk.Commit
	result, err1 := r.BuildAndSend([]sdk.Msg{msg}, request.BaseTx)
	if !err1.IsNil() {
		return "", err1
	}
	if !result.IsSuccess() {
		return "", sdk.Wrapf(result.GetLog())
	}

	requestID := result.GetTags().GetValue(TagRequestID)
	if request.Callback != nil {
		var subscription sdk.Subscription
		//TODO add query ?
		subscription, err = r.SubscribeNewBlockWithQuery(nil, func(block sdk.EventDataNewBlock) {
			tags := block.ResultBeginBlock.Tags
			r.Debug().
				Int64("height", block.Block.Height).
				Str("tags", tags.String()).
				Msg("received block")
			requestIDs := tags.GetValues(TagRequestID)
			for _, reqID := range requestIDs {
				if reqID == requestID {
					result, err := r.QueryRandom(requestID)
					var randomNum string
					if err.IsNil() {
						randomNum = result.RandomNum
					}
					request.Callback(requestID, randomNum, err)
					_ = r.Unscribe(subscription)
					return
				}
			}
		})
	}
	return requestID, sdk.Nil
}

// QueryRandom returns the random information of the specified reqID
func (r randomClient) QueryRandom(reqID string) (rpc.RandomInfo, sdk.Error) {
	param := struct {
		ReqID string
	}{
		ReqID: reqID,
	}

	var rand rand
	if err := r.QueryWithResponse("custom/rand/rand", param, &rand); err != nil {
		return rpc.RandomInfo{}, sdk.Wrap(err)
	}
	return rand.Convert().(rpc.RandomInfo), sdk.Nil
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
	return rs.Convert().([]rpc.RequestRandom), sdk.Nil
}
