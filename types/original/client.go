package original

import (
	"github.com/irisnet/irishub-sdk-go/utils/log"
	"github.com/tendermint/tendermint/libs/bytes"
)

type TxManager interface {
	BuildAndSend(msg []Msg, baseTx BaseTx) (ResultTx, Error)
	SendMsgBatch(msgs Msgs, baseTx BaseTx) ([]ResultTx, Error)
	Broadcast(signedTx StdTx, mode BroadcastMode) (ResultTx, Error)
}

type Queries interface {
	StoreQuery
	AccountQuery
	TxQuery
	ParamQuery
}

type ParamQuery interface {
	QueryParams(module string, res Response) Error
}

type StoreQuery interface {
	QueryWithResponse(path string, data interface{}, result Response) error
	Query(path string, data interface{}) ([]byte, error)
	QueryStore(key bytes.HexBytes, storeName string) (res []byte, err error)
}

type AccountQuery interface {
	QueryAccount(address string) (BaseAccount, Error)
	QueryAddress(name string) (AccAddress, Error)
}

type TxQuery interface {
	QueryTx(hash string) (ResultQueryTx, error)
	QueryTxs(builder *EventQueryBuilder, page, size int) (ResultSearchTxs, error)
}

type TokenManager interface {
	QueryToken(symbol string) (Token, error)
	SaveTokens(tokens ...Token)
}

type TokenConvert interface {
	ToMinCoin(coin ...DecCoin) (Coins, Error)
	ToMainCoin(coin ...Coin) (DecCoins, Error)
}

type Logger interface {
	Logger() *log.Logger
}

type BaseClient interface {
	TxManager
	TokenManager
	Queries
	TokenConvert
	TmClient
	Logger
}
