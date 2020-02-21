package tx

import "github.com/irisnet/irishub-sdk-go/types"

type Tx interface {
	BuildAndSend(msgs []types.Msg, baseTx types.BaseTx) (types.Result, error)
	Broadcast(signedTx types.StdTx, mode types.BroadcastMode) (types.Result, error)
	Sign(stdTx types.StdTx, name string, password string, offline bool) (types.StdTx, error)
}

type txClient struct {
	types.AbstractClient
}
