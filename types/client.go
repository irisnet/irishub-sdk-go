package types

import (
	"github.com/irisnet/irishub-sdk-go/tools/log"
	cmn "github.com/tendermint/tendermint/libs/common"
)

type TxManager interface {
	BuildAndSend(msg []Msg, baseTx BaseTx) (ResultTx, Error)
	SendMsgBatch(batch int, msgs []Msg, baseTx BaseTx) ([]ResultTx, Error)
	Broadcast(signedTx StdTx, mode BroadcastMode) (ResultTx, Error)
}

type Query interface {
	QueryWithResponse(path string, data interface{}, result Response) error
	Query(path string, data interface{}) ([]byte, error)
	QueryStore(key cmn.HexBytes, storeName string) (res []byte, err error)
	QueryAccount(address string) (BaseAccount, error)
	QueryAddress(name string) (addr AccAddress, err error)
	QueryToken(symbol string) (Token, error)
}

type TokenConvert interface {
	ToMinCoin(coin ...DecCoin) (Coins, error)
	ToMainCoin(coin ...Coin) (DecCoins, error)
}

type Logger interface {
	Logger() *log.Logger
}

type AbstractClient interface {
	TxManager
	Query
	TokenConvert
	TmClient
	Logger
}
