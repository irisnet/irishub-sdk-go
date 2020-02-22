package types

import (
	cmn "github.com/tendermint/tendermint/libs/common"
	tmclient "github.com/tendermint/tendermint/rpc/client"
)

type Queries interface {
	Query(path string, data cmn.HexBytes) ([]byte, error)
}

type WSClient interface {
	SubscribeNewBlock(callback EventNewBlockCallback) (Subscription, error)
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

type AbstractClient interface {
	Broadcast(baseTx BaseTx, msg []Msg) (Result, error)
	BroadcastTx(signedTx StdTx, mode BroadcastMode) (Result, error)
	Sign(stdTx StdTx, name string, password string, online bool) (StdTx, error)
	Query(path string, data interface{}, result interface{}) error
	QueryStore(key cmn.HexBytes, storeName string) ([]byte, error)
	QueryAccount(address string) (BaseAccount, error)
	QueryAddress(name string) AccAddress
}
