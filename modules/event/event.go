package event

import (
	"context"
	"encoding/base64"
	"fmt"
	"os/user"

	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Event interface {
	UnsubscribeAll() error
	SubscribeNewBlock(callback types.EventNewBlockCallback) (types.Subscription, error)
	SubscribeTx(builder *types.EventQueryBuilder, callback types.EventTxCallback) (types.Subscription, error)
}

type eventsClient struct {
	ws  *net.RPCClient
	cdc types.Codec
}

func New(ac types.AbstractClient) Event {
	ws := ac.GetRPC()
	return eventsClient{
		ws:  &ws,
		cdc: ac.GetCodec(),
	}
}

func (e eventsClient) start() {
	if !e.ws.IsRunning() {
		_ = e.ws.Start()
	}
}

func (e eventsClient) SubscribeNewBlock(callback types.EventNewBlockCallback) (types.Subscription, error) {
	ctx := context.Background()
	subscriber := getSubscriberID()
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlock).String()
	e.start()
	ch, err := e.ws.Subscribe(ctx, subscriber, query, 0)
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
				if err := e.cdc.UnmarshalBinaryLengthPrefixed(tx, &stdTx); err == nil {
					txs = append(txs, stdTx)
				}
			}

			var validatorUpdates []types.ValidatorUpdate
			for _, validator := range block.ResultEndBlock.ValidatorUpdates {
				var pubKey = types.EventPubKey{
					Type:  validator.PubKey.Type,
					Value: base64.StdEncoding.EncodeToString(validator.PubKey.Data),
				}
				validatorUpdates = append(validatorUpdates, types.ValidatorUpdate{
					PubKey: pubKey,
					Power:  validator.Power,
				})
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
					Tags: parseTags(block.ResultBeginBlock.Tags),
				},
				ResultEndBlock: types.ResultEndBlock{
					Tags:             parseTags(block.ResultEndBlock.Tags),
					ValidatorUpdates: validatorUpdates,
				},
			})
		}
	}()
	return types.NewSubscription(ctx, e.ws, query, subscriber), nil
}

func (e eventsClient) SubscribeTx(builder *types.EventQueryBuilder, callback types.EventTxCallback) (types.Subscription, error) {
	ctx := context.Background()
	subscriber := getSubscriberID()
	query := builder.Build()
	e.start()
	ch, err := e.ws.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return types.Subscription{}, err
	}
	go func() {
		for {
			data := <-ch
			tx := data.Data.(tmtypes.EventDataTx)
			var stdTx types.StdTx
			if err := e.cdc.UnmarshalBinaryLengthPrefixed(tx.Tx, &stdTx); err != nil {
				return
			}
			hash := cmn.HexBytes(tx.Tx.Hash()).String()
			result := types.TxResult{
				Log:       tx.Result.Log,
				GasWanted: tx.Result.GasWanted,
				GasUsed:   tx.Result.GasUsed,
				Tags:      parseTags(tx.Result.Tags),
			}
			dataTx := types.EventDataTx{
				Hash:   hash,
				Height: tx.Height,
				Index:  tx.Index,
				Tx:     stdTx,
				Result: result,
			}
			_ = e.ws.Unsubscribe(ctx, subscriber, query)
			callback(dataTx)
		}
	}()
	return types.NewSubscription(ctx, e.ws, query, subscriber), nil
}

func (e eventsClient) UnsubscribeAll() error {
	if e.ws.IsRunning() {
		return e.ws.UnsubscribeAll(context.Background(), getSubscriberID())
	}
	return nil
}

func getSubscriberID() string {
	u, err := user.Current()
	if err != nil {
		return "IRISHUB-SDK"
	}
	return fmt.Sprintf("subscriber-%s", u.Uid)
}

func parseTags(pairs []cmn.KVPair) (tags []types.Tag) {
	if pairs == nil || len(pairs) == 0 {
		return tags
	}
	for _, pair := range pairs {
		key := string(pair.Key)
		value := string(pair.Value)
		tags = append(tags, types.Tag{
			Key:   key,
			Value: value,
		})
	}
	return
}
