package modules

import (
	"context"
	"fmt"
	"os/user"

	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/irisnet/irishub-sdk-go/net"

	"github.com/irisnet/irishub-sdk-go/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type eventsClient struct {
	wsClient net.RPCClient
	cdc      types.Codec
}

type Subscriber struct {
	wsClient net.RPCClient
	ctx      context.Context
	query    string
	data     types.EventData
}

func (s Subscriber) Unsubscribe() {
	_ = s.wsClient.Unsubscribe(s.ctx, getSubscriberID(), s.query)
}
func (s Subscriber) GetData() types.EventData {
	return s.data
}

func NewEvent(tm types.TxManager) types.Event {
	wsClient := tm.GetTxContext().RPC
	_ = wsClient.Start()
	return eventsClient{
		wsClient: wsClient,
		cdc:      tm.GetTxContext().Codec,
	}
}

func (e eventsClient) SubscribeTx(query string, callback types.EventCallback) error {
	ctx := context.Background()
	subscriber := getSubscriberID()
	ch, err := e.wsClient.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return err
	}

	sub := Subscriber{
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
	sub := Subscriber{
		wsClient: e.wsClient,
		ctx:      ctx,
		query:    query,
	}
	go func() {
		for {
			data := <-ch
			block := data.Data.(tmtypes.EventDataNewBlock)
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
