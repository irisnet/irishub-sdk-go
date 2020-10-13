package original

import (
	"context"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/irisnet/irishub-sdk-go/utils/log"
	"github.com/irisnet/irishub-sdk-go/utils/uuid"
	"github.com/tendermint/tendermint/libs/bytes"
	rpc "github.com/tendermint/tendermint/rpc/client"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
)

type rpcClient struct {
	rpc.Client
	*log.Logger
	cdc original.Codec
}

func NewRPCClient(remote string, cdc original.Codec, log *log.Logger) original.TmClient {
	client, err := rpchttp.New(remote, "/websocket")
	if err != nil {
		panic(err)
	}
	_ = client.Start()
	return rpcClient{
		Client: client,
		Logger: log,
		cdc:    cdc,
	}
}

//=============================================================================
//SubscribeNewBlock implement WSClient interface
func (r rpcClient) SubscribeNewBlock(builder *original.EventQueryBuilder,
	handler original.EventNewBlockHandler) (original.Subscription, original.Error) {
	if builder == nil {
		builder = original.NewEventQueryBuilder()
	}

	builder.AddCondition(original.Cond(original.TypeKey).EQ(tmtypes.EventNewBlock))
	query := builder.Build()

	return r.SubscribeAny(query, func(data original.EventData) {
		handler(data.(original.EventDataNewBlock))
	})
}

//SubscribeTx implement WSClient interface
func (r rpcClient) SubscribeTx(builder *original.EventQueryBuilder, handler original.EventTxHandler) (original.Subscription, original.Error) {
	if builder == nil {
		builder = original.NewEventQueryBuilder()
	}
	query := builder.AddCondition(original.Cond(original.TypeKey).EQ(original.TxValue)).Build()
	return r.SubscribeAny(query, func(data original.EventData) {
		handler(data.(original.EventDataTx))
	})
}

func (r rpcClient) SubscribeNewBlockHeader(handler original.EventNewBlockHeaderHandler) (original.Subscription, original.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	return r.SubscribeAny(query, func(data original.EventData) {
		handler(data.(original.EventDataNewBlockHeader))
	})
}

func (r rpcClient) SubscribeValidatorSetUpdates(handler original.EventValidatorSetUpdatesHandler) (original.Subscription, original.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventValidatorSetUpdates).String()
	return r.SubscribeAny(query, func(data original.EventData) {
		handler(data.(original.EventDataValidatorSetUpdates))
	})
}

func (r rpcClient) Resubscribe(subscription original.Subscription, handler original.EventHandler) (err original.Error) {
	_, err = r.SubscribeAny(subscription.Query, handler)
	return
}

func (r rpcClient) Unsubscribe(subscription original.Subscription) original.Error {
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
		return original.Wrap(err)
	}
	return nil
}

func (r rpcClient) SubscribeAny(query string, handler original.EventHandler) (subscription original.Subscription, err original.Error) {
	ctx := context.Background()
	subscriber := getSubscriber()
	ch, e := r.Subscribe(ctx, subscriber, query, 0)
	if e != nil {
		return subscription, original.Wrap(e)
	}

	r.Info().
		Str("query", query).
		Str("subscriber", subscriber).
		Msg("subscribe event")

	subscription = original.Subscription{
		Ctx:   ctx,
		Query: query,
		ID:    subscriber,
	}

	go func() {
		for {
			data := <-ch
			go func() {
				defer original.CatchPanic(func(errMsg string) {
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

func (r rpcClient) parseTx(data original.EventData) original.EventDataTx {
	tx := data.(tmtypes.EventDataTx)
	var stdTx original.StdTx
	if err := r.cdc.UnmarshalBinaryBare(tx.Tx, &stdTx); err != nil {
		return original.EventDataTx{}
	}
	hash := bytes.HexBytes(tx.Tx.Hash()).String()
	result := original.TxResult{
		Code:      tx.Result.Code,
		Log:       tx.Result.Log,
		GasWanted: tx.Result.GasWanted,
		GasUsed:   tx.Result.GasUsed,
		//Tags:      sdk.ParseTags(tx.Result.Tags),
	}
	return original.EventDataTx{
		Hash:   hash,
		Height: tx.Height,
		Index:  tx.Index,
		Tx:     stdTx,
		Result: result,
	}
}

func (r rpcClient) parseNewBlock(data original.EventData) original.EventDataNewBlock {
	block := data.(tmtypes.EventDataNewBlock)
	return original.EventDataNewBlock{
		Block:            original.ParseBlock(r.cdc, block.Block),
		ResultBeginBlock: original.ResultBeginBlock{
			//Tags: sdk.ParseTags(block.ResultBeginBlock.Tags),
		},
		ResultEndBlock: original.ResultEndBlock{
			//Tags:             sdk.ParseTags(block.ResultEndBlock.Tags),
			ValidatorUpdates: original.ParseValidatorUpdate(block.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseNewBlockHeader(data original.EventData) original.EventDataNewBlockHeader {
	blockHeader := data.(tmtypes.EventDataNewBlockHeader)
	return original.EventDataNewBlockHeader{
		Header:           blockHeader.Header,
		ResultBeginBlock: original.ResultBeginBlock{},
		ResultEndBlock: original.ResultEndBlock{
			ValidatorUpdates: original.ParseValidatorUpdate(blockHeader.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseValidatorSetUpdates(data original.EventData) original.EventDataValidatorSetUpdates {
	validatorSet := data.(tmtypes.EventDataValidatorSetUpdates)
	return original.EventDataValidatorSetUpdates{
		ValidatorUpdates: original.ParseValidators(validatorSet.ValidatorUpdates),
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
