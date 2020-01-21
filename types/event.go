package types

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"
	tmquery "github.com/tendermint/tendermint/libs/pubsub/query"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Event interface {
	UnsubscribeAll() error
	SubscribeNewBlock(callback EventCallback) error
	SubscribeTx(query string, callback EventCallback) error
}

type Subscription interface {
	Unsubscribe()
	GetData() EventData
}

type EventCallback func(Subscription)
type EventData interface{}
type EventDataTx struct {
	Hash   string                 `json:"hash"`
	Height int64                  `json:"height"`
	Index  uint32                 `json:"index"`
	Tx     StdTx                  `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
}

func EventQueryTxFor(txHash string) string {
	query := tmquery.MustParse(fmt.Sprintf("%s='%s' AND %s='%s'",
		tmtypes.EventTypeKey, tmtypes.EventTx, tmtypes.TxHashKey, txHash))
	return query.String()
}
