package types

import (
	cmn "github.com/tendermint/tendermint/libs/common"
	tmclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/irisnet/irishub-sdk-go/tools/log"
)

type Queries interface {
	Query(path string, data cmn.HexBytes) ([]byte, error)
}

type WSClient interface {
	SubscribeNewBlock(callback EventNewBlockCallback) (Subscription, error)
	SubscribeNewBlockWithQuery(builder *EventQueryBuilder, callback EventNewBlockCallback) (Subscription, error)
	SubscribeTx(builder *EventQueryBuilder, callback EventTxCallback) (Subscription, error)
	SubscribeNewBlockHeader(callback EventNewBlockHeaderCallback) (Subscription, error)
	SubscribeValidatorSetUpdates(callback EventValidatorSetUpdatesCallback) (Subscription, error)
	Unscribe(subscription Subscription) error
}

type TmClient interface {
	tmclient.Client
	WSClient
	Queries
}

type TxManager interface {
	BuildAndSend(msg []Msg, baseTx BaseTx) (Result, Error)
	Broadcast(signedTx StdTx, mode BroadcastMode) (Result, Error)
}

type Query interface {
	QueryWithResponse(path string, data interface{}, result Response) error
	Query(path string, data interface{}) ([]byte, error)
	QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error)
	QueryAccount(address string) (BaseAccount, error)
	QueryAddress(name, password string) (addr AccAddress, err error)
}

type Logger interface {
	Logger() *log.Logger
}

type AbstractClient interface {
	TxManager
	Query
	WSClient
	Logger
}
