package net

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/user"

	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpc "github.com/tendermint/tendermint/rpc/client"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/irisnet/irishub-sdk-go/types"
)

type RPCClient struct {
	rpc.Client
	cdc types.Codec
}

func NewRPCClient(remote string, cdc types.Codec) RPCClient {
	client := rpc.NewHTTP(remote, "/websocket")
	return RPCClient{
		Client: client,
		cdc:    cdc,
	}
}

func (r RPCClient) Query(path string, data cmn.HexBytes) (res []byte, err error) {
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

func (r RPCClient) start() {
	if !r.IsRunning() {
		_ = r.Start()
	}
}

//=============================================================================

//SubscribeNewBlock implement WSClient interface
func (r RPCClient) SubscribeNewBlock(callback types.EventNewBlockCallback) (types.Subscription, error) {
	return r.SubscribeNewBlockWithParams(nil, callback)
}

//SubscribeNewBlock implement WSClient interface
func (r RPCClient) SubscribeNewBlockWithParams(builder *types.EventQueryBuilder, callback types.EventNewBlockCallback) (types.Subscription, error) {
	ctx := context.Background()
	subscriber := getSubscriberID()
	if builder == nil {
		builder = types.NewEventQueryBuilder()
	}
	builder.AddCondition(types.TypeKey, tmtypes.EventNewBlock)
	query := builder.Build()
	r.start()
	ch, err := r.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return types.Subscription{}, nil
	}
	go func() {
		for {
			data := <-ch
			block := data.Data.(tmtypes.EventDataNewBlock)
			var txs []types.StdTx
			for _, tx := range block.Block.Data.Txs {
				var stdTx types.StdTx
				if err := r.cdc.UnmarshalBinaryLengthPrefixed(tx, &stdTx); err == nil {
					txs = append(txs, stdTx)
				}
			}

			callback(types.EventDataNewBlock{
				Block: types.Block{
					Header: block.Block.Header,
					Data: types.Data{
						Txs: txs,
					},
					Evidence:   block.Block.Evidence,
					LastCommit: block.Block.LastCommit,
				},
				ResultBeginBlock: types.ResultBeginBlock{
					Tags: types.ParseTags(block.ResultBeginBlock.Tags),
				},
				ResultEndBlock: types.ResultEndBlock{
					Tags:             types.ParseTags(block.ResultEndBlock.Tags),
					ValidatorUpdates: parseValidatorUpdate(block.ResultEndBlock.ValidatorUpdates),
				},
			})
		}
	}()
	return types.NewSubscription(ctx, query, subscriber), nil
}

//SubscribeTx implement WSClient interface
func (r RPCClient) SubscribeTx(builder *types.EventQueryBuilder, callback types.EventTxCallback) (types.Subscription, error) {
	ctx := context.Background()
	subscriber := getSubscriberID()
	query := builder.AddCondition(types.TypeKey, types.TxValue).Build()
	r.start()
	ch, err := r.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return types.Subscription{}, err
	}

	go func() {
		for {
			data := <-ch
			tx := data.Data.(tmtypes.EventDataTx)
			var stdTx types.StdTx
			if err := r.cdc.UnmarshalBinaryLengthPrefixed(tx.Tx, &stdTx); err != nil {
				return
			}
			hash := cmn.HexBytes(tx.Tx.Hash()).String()
			result := types.TxResult{
				Log:       tx.Result.Log,
				GasWanted: tx.Result.GasWanted,
				GasUsed:   tx.Result.GasUsed,
				Tags:      types.ParseTags(tx.Result.Tags),
			}
			dataTx := types.EventDataTx{
				Hash:   hash,
				Height: tx.Height,
				Index:  tx.Index,
				Tx:     stdTx,
				Result: result,
			}
			callback(dataTx)
		}
	}()
	return types.NewSubscription(ctx, query, subscriber), nil
}

func (r RPCClient) SubscribeNewBlockHeader(callback types.EventNewBlockHeaderCallback) (types.Subscription, error) {
	ctx := context.Background()
	subscriber := getSubscriberID()
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlockHeader).String()
	r.start()
	ch, err := r.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return types.Subscription{}, nil
	}

	go func() {
		for {
			data := <-ch
			blockHeader := data.Data.(tmtypes.EventDataNewBlockHeader)
			callback(types.EventDataNewBlockHeader{
				Header: blockHeader.Header,
				ResultBeginBlock: types.ResultBeginBlock{
					Tags: types.ParseTags(blockHeader.ResultBeginBlock.Tags),
				},
				ResultEndBlock: types.ResultEndBlock{
					Tags:             types.ParseTags(blockHeader.ResultEndBlock.Tags),
					ValidatorUpdates: parseValidatorUpdate(blockHeader.ResultEndBlock.ValidatorUpdates),
				},
			})
		}
	}()
	return types.NewSubscription(ctx, query, subscriber), nil
}

func (r RPCClient) SubscribeValidatorSetUpdates(callback types.EventValidatorSetUpdatesCallback) (types.Subscription, error) {
	ctx := context.Background()
	subscriber := getSubscriberID()
	query := tmtypes.QueryForEvent(tmtypes.EventValidatorSetUpdates).String()
	r.start()
	ch, err := r.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return types.Subscription{}, nil
	}

	go func() {
		for {
			data := <-ch
			validatorSet := data.Data.(tmtypes.EventDataValidatorSetUpdates)
			callback(types.EventDataValidatorSetUpdates{
				ValidatorUpdates: parseValidators(validatorSet.ValidatorUpdates),
			})
		}
	}()
	return types.NewSubscription(ctx, query, subscriber), nil
}

func (r RPCClient) Unscribe(subscription types.Subscription) error {
	return r.Client.Unsubscribe(subscription.Ctx, subscription.ID, subscription.Query)
}

func getSubscriberID() string {
	u, err := user.Current()
	if err != nil {
		return "IRISHUB-SDK"
	}
	return fmt.Sprintf("subscriber-%s", u.Uid)
}

func parseValidatorUpdate(vp abcitypes.ValidatorUpdates) (validatorUpdates []types.ValidatorUpdate) {
	for _, validator := range vp {
		var pubKey = types.EventPubKey{
			Type:  validator.PubKey.Type,
			Value: base64.StdEncoding.EncodeToString(validator.PubKey.Data),
		}
		validatorUpdates = append(validatorUpdates, types.ValidatorUpdate{
			PubKey: pubKey,
			Power:  validator.Power,
		})
	}
	return
}

func parseValidators(valSet []*tmtypes.Validator) (validators []types.Validator) {
	for _, v := range valSet {
		valAddr, _ := types.ConsAddressFromHex(v.Address.String())
		pubKey, _ := types.Bech32ifyConsPub(v.PubKey)
		validators = append(validators, types.Validator{
			Address:          valAddr.String(),
			PubKey:           pubKey,
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		})
	}
	return
}
