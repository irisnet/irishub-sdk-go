package client

import (
	"context"
	"fmt"

	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/irisnet/irishub-sdk-go/tools/log"
	"github.com/irisnet/irishub-sdk-go/tools/uuid"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type rpcClient struct {
	rpc.Client
	*log.Logger
	cdc sdk.Codec
}

func NewRPCClient(remote string, cdc sdk.Codec, log *log.Logger) sdk.TmClient {
	client := rpc.NewHTTP(remote, "/websocket")
	_ = client.Start()
	return rpcClient{
		Client: client,
		Logger: log,
		cdc:    cdc,
	}
}

//=============================================================================
//SubscribeNewBlock implement WSClient interface
func (r rpcClient) SubscribeNewBlock(builder *sdk.EventQueryBuilder,
	handler sdk.EventNewBlockHandler) (sdk.Subscription, sdk.Error) {
	if builder == nil {
		builder = sdk.NewEventQueryBuilder()
	}

	builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(tmtypes.EventNewBlock))
	query := builder.Build()

	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataNewBlock))
	})
}

//SubscribeTx implement WSClient interface
func (r rpcClient) SubscribeTx(builder *sdk.EventQueryBuilder, handler sdk.EventTxHandler) (sdk.Subscription, sdk.Error) {
	if builder == nil {
		builder = sdk.NewEventQueryBuilder()
	}
	query := builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(sdk.TxValue)).Build()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataTx))
	})
}

func (r rpcClient) SubscribeNewBlockHeader(handler sdk.EventNewBlockHeaderHandler) (sdk.Subscription, sdk.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataNewBlockHeader))
	})
}

func (r rpcClient) SubscribeValidatorSetUpdates(handler sdk.EventValidatorSetUpdatesHandler) (sdk.Subscription, sdk.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventValidatorSetUpdates).String()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		handler(data.(sdk.EventDataValidatorSetUpdates))
	})
}

func (r rpcClient) Resubscribe(subscription sdk.Subscription, handler sdk.EventHandler) (err sdk.Error) {
	_, err = r.SubscribeAny(subscription.Query, handler)
	return
}

func (r rpcClient) Unsubscribe(subscription sdk.Subscription) sdk.Error {
	r.Info().
		Str("query", subscription.Query).
		Str("subscriber", subscription.ID).
		Msg("end to subscribe event")
	err := r.Client.Unsubscribe(subscription.Ctx, subscription.ID, subscription.Query)
	if err != nil {
		r.Err(err).
			Str("query", subscription.Query).
			Str("subscriber", subscription.ID).
			Msg("unsubscribe failed")
		return sdk.Wrap(err)
	}
	return nil
}

func (r rpcClient) SubscribeAny(query string, handler sdk.EventHandler) (subscription sdk.Subscription, err sdk.Error) {
	ctx := context.Background()
	subscriber := getSubscriber()
	ch, e := r.Subscribe(ctx, subscriber, query, 0)
	if e != nil {
		return subscription, sdk.Wrap(e)
	}

	r.Info().
		Str("query", query).
		Str("subscriber", subscriber).
		Msg("subscribe event")

	subscription = sdk.Subscription{
		Ctx:   ctx,
		Query: query,
		ID:    subscriber,
	}

	go func() {
		for {
			data := <-ch
			go func() {
				defer sdk.CatchPanic(func(errMsg string) {
					r.Error().
						Str("query", query).
						Str("subscriber", subscriber).
						Msgf("subscribe event failed:%s", errMsg)
				})

				switch data := data.Data.(type) {
				case tmtypes.EventDataTx:
					handler(r.parseTx(data))
					return
				case tmtypes.EventDataNewBlock:
					handler(r.parseNewBlock(data))
					return
				case tmtypes.EventDataNewBlockHeader:
					handler(r.parseNewBlockHeader(data))
					return
				case tmtypes.EventDataValidatorSetUpdates:
					handler(r.parseValidatorSetUpdates(data))
					return
				default:
					handler(data)
				}
			}()
		}
	}()
	return
}

func (r rpcClient) parseTx(data sdk.EventData) sdk.EventDataTx {
	tx := data.(tmtypes.EventDataTx)
	var stdTx sdk.StdTx
	if err := r.cdc.UnmarshalBinaryLengthPrefixed(tx.Tx, &stdTx); err != nil {
		return sdk.EventDataTx{}
	}
	hash := cmn.HexBytes(tx.Tx.Hash()).String()
	result := sdk.TxResult{
		Code:      tx.Result.Code,
		Log:       tx.Result.Log,
		GasWanted: tx.Result.GasWanted,
		GasUsed:   tx.Result.GasUsed,
		Tags:      sdk.ParseTags(tx.Result.Tags),
	}
	return sdk.EventDataTx{
		Hash:   hash,
		Height: tx.Height,
		Index:  tx.Index,
		Tx:     stdTx,
		Result: result,
	}
}

func (r rpcClient) parseNewBlock(data sdk.EventData) sdk.EventDataNewBlock {
	block := data.(tmtypes.EventDataNewBlock)
	return sdk.EventDataNewBlock{
		Block: sdk.ParseBlock(r.cdc, block.Block),
		ResultBeginBlock: sdk.ResultBeginBlock{
			Tags: sdk.ParseTags(block.ResultBeginBlock.Tags),
		},
		ResultEndBlock: sdk.ResultEndBlock{
			Tags:             sdk.ParseTags(block.ResultEndBlock.Tags),
			ValidatorUpdates: sdk.ParseValidatorUpdate(block.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseNewBlockHeader(data sdk.EventData) sdk.EventDataNewBlockHeader {
	blockHeader := data.(tmtypes.EventDataNewBlockHeader)
	return sdk.EventDataNewBlockHeader{
		Header: blockHeader.Header,
		ResultBeginBlock: sdk.ResultBeginBlock{
			Tags: sdk.ParseTags(blockHeader.ResultBeginBlock.Tags),
		},
		ResultEndBlock: sdk.ResultEndBlock{
			Tags:             sdk.ParseTags(blockHeader.ResultEndBlock.Tags),
			ValidatorUpdates: sdk.ParseValidatorUpdate(blockHeader.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseValidatorSetUpdates(data sdk.EventData) sdk.EventDataValidatorSetUpdates {
	validatorSet := data.(tmtypes.EventDataValidatorSetUpdates)
	return sdk.EventDataValidatorSetUpdates{
		ValidatorUpdates: sdk.ParseValidators(validatorSet.ValidatorUpdates),
	}
}

func getSubscriber() string {
	subscriber := "irishub-sdk-go"
	id, err := uuid.NewV1()
	if err == nil {
		subscriber = fmt.Sprintf("%s-%s", subscriber, id.String())
	}
	return subscriber
}
