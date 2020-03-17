package client

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/tools/log"

	"github.com/irisnet/irishub-sdk-go/tools/uuid"
	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type rpcClient struct {
	rpc.Client
	cdc sdk.Codec
	*log.Logger
}

func NewRPCClient(remote string, cdc sdk.Codec, log *log.Logger) sdk.TmClient {
	client := rpc.NewHTTP(remote, "/websocket")
	return rpcClient{
		Client: client,
		cdc:    cdc,
		Logger: log,
	}
}

func (r rpcClient) Query(path string, data cmn.HexBytes) (res []byte, err error) {
	result, err := r.ABCIQueryWithOptions(path, data, rpc.DefaultABCIQueryOptions)
	if err != nil {
		return res, err
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, errors.Errorf(resp.Log)
	}

	return resp.Value, nil
}

func (r rpcClient) start() {
	if !r.IsRunning() {
		_ = r.Start()
	}
}

//=============================================================================

//SubscribeNewBlock implement WSClient interface
func (r rpcClient) SubscribeNewBlock(callback sdk.EventNewBlockCallback) (sdk.Subscription, error) {
	return r.SubscribeNewBlockWithQuery(nil, callback)
}

//SubscribeNewBlock implement WSClient interface
func (r rpcClient) SubscribeNewBlockWithQuery(builder *sdk.EventQueryBuilder, callback sdk.EventNewBlockCallback) (sdk.Subscription, error) {
	if builder == nil {
		builder = sdk.NewEventQueryBuilder()
	}

	builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(tmtypes.EventNewBlock))
	query := builder.Build()

	return r.subscribe(query, func(data tmtypes.TMEventData) {
		block := data.(tmtypes.EventDataNewBlock)
		var txs []sdk.StdTx
		for _, tx := range block.Block.Data.Txs {
			var stdTx sdk.StdTx
			if err := r.cdc.UnmarshalBinaryLengthPrefixed(tx, &stdTx); err == nil {
				txs = append(txs, stdTx)
			}
		}

		callback(sdk.EventDataNewBlock{
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
		})
	})
}

//SubscribeTx implement WSClient interface
func (r rpcClient) SubscribeTx(builder *sdk.EventQueryBuilder, callback sdk.EventTxCallback) (sdk.Subscription, error) {
	query := builder.AddCondition(sdk.Cond(sdk.TypeKey).EQ(sdk.TxValue)).Build()
	return r.subscribe(query, func(data tmtypes.TMEventData) {
		tx := data.(tmtypes.EventDataTx)
		var stdTx sdk.StdTx
		if err := r.cdc.UnmarshalBinaryLengthPrefixed(tx.Tx, &stdTx); err != nil {
			return
		}
		hash := cmn.HexBytes(tx.Tx.Hash()).String()
		result := sdk.TxResult{
			Log:       tx.Result.Log,
			GasWanted: tx.Result.GasWanted,
			GasUsed:   tx.Result.GasUsed,
			Tags:      sdk.ParseTags(tx.Result.Tags),
		}
		dataTx := sdk.EventDataTx{
			Hash:   hash,
			Height: tx.Height,
			Index:  tx.Index,
			Tx:     stdTx,
			Result: result,
		}
		callback(dataTx)
	})
}

func (r rpcClient) SubscribeNewBlockHeader(callback sdk.EventNewBlockHeaderCallback) (sdk.Subscription, error) {
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	return r.subscribe(query, func(data tmtypes.TMEventData) {
		blockHeader := data.(tmtypes.EventDataNewBlockHeader)
		callback(sdk.EventDataNewBlockHeader{
			Header: blockHeader.Header,
			ResultBeginBlock: sdk.ResultBeginBlock{
				Tags: sdk.ParseTags(blockHeader.ResultBeginBlock.Tags),
			},
			ResultEndBlock: sdk.ResultEndBlock{
				Tags:             sdk.ParseTags(blockHeader.ResultEndBlock.Tags),
				ValidatorUpdates: parseValidatorUpdate(blockHeader.ResultEndBlock.ValidatorUpdates),
			},
		})
	})
}

func (r rpcClient) SubscribeValidatorSetUpdates(callback sdk.EventValidatorSetUpdatesCallback) (sdk.Subscription, error) {
	query := tmtypes.QueryForEvent(tmtypes.EventValidatorSetUpdates).String()
	return r.subscribe(query, func(data tmtypes.TMEventData) {
		validatorSet := data.(tmtypes.EventDataValidatorSetUpdates)
		callback(sdk.EventDataValidatorSetUpdates{
			ValidatorUpdates: parseValidators(validatorSet.ValidatorUpdates),
		})
	})
}

func (r rpcClient) subscribe(query string, callback func(data tmtypes.TMEventData)) (sdk.Subscription, error) {
	ctx := context.Background()
	subscriber := getSubscriber()
	r.start()

	ch, err := r.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return sdk.Subscription{}, nil
	}

	r.Info().
		Str("query", query).
		Str("subscriber", subscriber).
		Msg("begin to subscribe event")

	go func() {
		for {
			data := <-ch
			callback(data.Data)
		}
	}()
	return sdk.NewSubscription(ctx, query, subscriber), nil
}

func (r rpcClient) Unscribe(subscription sdk.Subscription) error {
	r.Info().
		Str("query", subscription.Query).
		Str("subscriber", subscription.ID).
		Msg("end to subscribe event")
	return r.Client.Unsubscribe(subscription.Ctx, subscription.ID, subscription.Query)
}

func getSubscriber() string {
	subscriber := "irishub-sdk-go"
	id, err := uuid.NewV1()
	if err == nil {
		subscriber = fmt.Sprintf("%s-%s", subscriber, id.String())
	}
	return subscriber
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

func parseValidators(valSet []*tmtypes.Validator) (validators []sdk.Validator) {
	for _, v := range valSet {
		valAddr, _ := sdk.ConsAddressFromHex(v.Address.String())
		pubKey, _ := sdk.Bech32ifyConsPub(v.PubKey)
		validators = append(validators, sdk.Validator{
			Address:          valAddr.String(),
			PubKey:           pubKey,
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		})
	}
	return
}
