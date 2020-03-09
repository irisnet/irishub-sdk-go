package random

import (
	"errors"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type randomClient struct {
	sdk.AbstractClient
	*log.Logger
}

func New(ac sdk.AbstractClient) rpc.Random {
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
func (r randomClient) Generate(request rpc.RandomRequest) (string, error) {
	consumer, err := r.QueryAddress(request.From, request.Password)
	if err != nil {
		return "", err
	}

	msg := MsgRequestRand{
		Consumer:      consumer,
		BlockInterval: request.BlockInterval,
	}

	//mode must be set to commit
	request.BaseTx.Mode = sdk.Commit
	result, err := r.Broadcast(request.BaseTx, []sdk.Msg{msg})
	if err != nil {
		return "", err
	}
	if !result.IsSuccess() {
		return "", errors.New(result.GetLog())
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
					if err == nil {
						randomNum = result.RandomNum
						r.Debug().
							Int64("height", block.Block.Height).
							Str("requestID", reqID).
							Str("txHash", result.RequestTxHash).
							Msg("received block")
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
func (r randomClient) QueryRandom(reqID string) (rpc.RandomInfo, error) {
	param := struct {
		ReqID string
	}{
		ReqID: reqID,
	}

	var rand rand
	if err := r.QueryWithResponse("custom/rand/rand", param, &rand); err != nil {
		return rpc.RandomInfo{}, err
	}
	return rand.Convert().(rpc.RandomInfo), nil
}

// QueryRequests returns the list of request by the specified block height
func (r randomClient) QueryRequests(height int64) ([]rpc.RequestRandom, error) {
	param := struct {
		Height int64
	}{
		Height: height,
	}

	var rs requests
	if err := r.QueryWithResponse("custom/rand/queue", param, &rs); err != nil {
		return nil, err
	}
	return rs.Convert().([]rpc.RequestRandom), nil
}
