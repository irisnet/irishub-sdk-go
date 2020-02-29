package types

import (
	"github.com/irisnet/irishub-sdk-go/tools/log"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

type Queries interface {
	Query(path string, data cmn.HexBytes) ([]byte, error)
}

type WSClient interface {
	SubscribeNewBlock(callback EventNewBlockCallback) (Subscription, error)
	SubscribeNewBlockWithParams(builder *EventQueryBuilder, callback EventNewBlockCallback) (Subscription, error)
	SubscribeTx(builder *EventQueryBuilder, callback EventTxCallback) (Subscription, error)
	SubscribeNewBlockHeader(callback EventNewBlockHeaderCallback) (Subscription, error)
	SubscribeValidatorSetUpdates(callback EventValidatorSetUpdatesCallback) (Subscription, error)
	Unscribe(subscription Subscription) error
}

type RPC interface {
	tmclient.Client
	WSClient
	Queries
}

type TxManager interface {
	Broadcast(baseTx BaseTx, msg []Msg) (Result, error)
	BroadcastTx(signedTx StdTx, mode BroadcastMode) (Result, error)
	Sign(stdTx StdTx, name string, password string, online bool) (StdTx, error)
}

type Query interface {
	Query(path string, data interface{}, result interface{}) error
	QueryAccount(address string) (BaseAccount, error)
	QueryAddress(name, password string) (addr AccAddress, err error)
}

type AbstractClient interface {
	TxManager
	Query
	WSClient
	Logger() *log.Logger
}
