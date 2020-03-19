package client

import (
	"context"
	"encoding/base64"
	"fmt"

	abcitypes "github.com/tendermint/tendermint/abci/types"

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
	callback sdk.EventNewBlockCallback) (sdk.Subscription, sdk.Error) {
	if builder == nil {
		builder = sdk.NewEventQueryBuilder()
	}

	builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(tmtypes.EventNewBlock))
	query := builder.Build()

	return r.SubscribeAny(query, func(data sdk.EventData) {
		callback(data.(sdk.EventDataNewBlock))
	})
}

//SubscribeTx implement WSClient interface
func (r rpcClient) SubscribeTx(builder *sdk.EventQueryBuilder, callback sdk.EventTxCallback) (sdk.Subscription, sdk.Error) {
	if builder == nil {
		builder = sdk.NewEventQueryBuilder()
	}
	query := builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(sdk.TxValue)).Build()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		callback(data.(sdk.EventDataTx))
	})
}

func (r rpcClient) SubscribeNewBlockHeader(callback sdk.EventNewBlockHeaderCallback) (sdk.Subscription, sdk.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		callback(data.(sdk.EventDataNewBlockHeader))
	})
}

func (r rpcClient) SubscribeValidatorSetUpdates(callback sdk.EventValidatorSetUpdatesCallback) (sdk.Subscription, sdk.Error) {
	query := tmtypes.QueryForEvent(tmtypes.EventValidatorSetUpdates).String()
	return r.SubscribeAny(query, func(data sdk.EventData) {
		callback(data.(sdk.EventDataValidatorSetUpdates))
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
	var txs []sdk.StdTx
	for _, tx := range block.Block.Data.Txs {
		var stdTx sdk.StdTx
		if err := r.cdc.UnmarshalBinaryLengthPrefixed(tx, &stdTx); err == nil {
			txs = append(txs, stdTx)
		}
	}

	return sdk.EventDataNewBlock{
		Block: sdk.Block{
			Header: block.Block.Header,
			Data: sdk.Data{
				Txs: txs,
			},
			Evidence:   block.Block.Evidence,
			LastCommit: block.Block.LastCommit,
		},
		ResultBeginBlock: sdk.ResultBeginBlock{
			Tags: sdk.ParseTags(block.ResultBeginBlock.Tags),
		},
		ResultEndBlock: sdk.ResultEndBlock{
			Tags:             sdk.ParseTags(block.ResultEndBlock.Tags),
			ValidatorUpdates: parseValidatorUpdate(block.ResultEndBlock.ValidatorUpdates),
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
			ValidatorUpdates: parseValidatorUpdate(blockHeader.ResultEndBlock.ValidatorUpdates),
		},
	}
}

func (r rpcClient) parseValidatorSetUpdates(data sdk.EventData) sdk.EventDataValidatorSetUpdates {
	validatorSet := data.(tmtypes.EventDataValidatorSetUpdates)

	var validators []sdk.Validator
	for _, v := range validatorSet.ValidatorUpdates {
		valAddr, _ := sdk.ConsAddressFromHex(v.Address.String())
		pubKey, _ := sdk.Bech32ifyConsPub(v.PubKey)
		validators = append(validators, sdk.Validator{
			Address:          valAddr.String(),
			PubKey:           pubKey,
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		})
	}
	return sdk.EventDataValidatorSetUpdates{
		ValidatorUpdates: validators,
	}
}

func parseValidatorUpdate(vp abcitypes.ValidatorUpdates) (validatorUpdates []sdk.ValidatorUpdate) {
	for _, validator := range vp {
		var pubKey = sdk.EventPubKey{
			Type:  validator.PubKey.Type,
			Value: base64.StdEncoding.EncodeToString(validator.PubKey.Data),
		}
		validatorUpdates = append(validatorUpdates, sdk.ValidatorUpdate{
			PubKey: pubKey,
			Power:  validator.Power,
		})
	}
	return
}

func getSubscriber() string {
	subscriber := "irishub-sdk-go"
	id, err := uuid.NewV1()
	if err == nil {
		subscriber = fmt.Sprintf("%s-%s", subscriber, id.String())
	}
	return subscriber
}
