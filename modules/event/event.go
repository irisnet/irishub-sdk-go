package event

import (
	"context"
	"fmt"
	"os/user"

	"github.com/irisnet/irishub-sdk-go/net"
	"github.com/irisnet/irishub-sdk-go/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Event interface {
	UnsubscribeAll() error
	SubscribeNewBlock(callback types.EventCallback) error
	SubscribeTx(params types.KVPair, callback types.EventCallback) error
}

type eventsClient struct {
	wsClient net.RPCClient
	cdc      types.Codec
}

type subContent struct {
	wsClient net.RPCClient
	ctx      context.Context
	query    string
	data     types.EventData
}

func (s subContent) Unsubscribe() {
	_ = s.wsClient.Unsubscribe(s.ctx, getSubscriberID(), s.query)
}
func (s subContent) GetData() types.EventData {
	return s.data
}

func NewEvent(tm types.TxCtxManager) Event {
	wsClient := tm.GetRPC()
	_ = wsClient.Start()
	return eventsClient{
		wsClient: wsClient,
		cdc:      tm.GetCodec(),
	}
}

func (e eventsClient) SubscribeTx(params types.KVPair, callback types.EventCallback) error {
	ctx := context.Background()
	subscriber := getSubscriberID()
	query := params.ToQueryString()
	ch, err := e.wsClient.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return err
	}

	sub := subContent{
		wsClient: e.wsClient,
		ctx:      ctx,
		query:    query,
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
			result := types.EventDataTx{
				Hash:   hash,
				Height: tx.Height,
				Index:  tx.Index,
				Tx:     stdTx,
				Result: tx.Result,
			}
			sub.data = result
			_ = e.wsClient.Unsubscribe(ctx, subscriber, query)
			callback(sub)
		}
	}()
	return nil
}

func (e eventsClient) SubscribeNewBlock(callback types.EventCallback) error {
	ctx := context.Background()
	subscriber := getSubscriberID()
	query := tmtypes.QueryForEvent(tmtypes.EventNewBlock).String()
	ch, err := e.wsClient.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return err
	}
	sub := subContent{
		wsClient: e.wsClient,
		ctx:      ctx,
		query:    query,
	}
	go func() {
		for {
			data := <-ch
			block := data.Data.(types.EventDataNewBlock)
			sub.data = block
			callback(sub)
		}
	}()
	return nil
}

func (e eventsClient) UnsubscribeAll() error {
	return e.wsClient.UnsubscribeAll(context.Background(), getSubscriberID())
}

func getSubscriberID() string {
	u, err := user.Current()
	if err != nil {
		return "IRISHUB-SDK"
	}
	return fmt.Sprintf("subscriber-%s", u.Uid)
}
