package tx

import "github.com/irisnet/irishub-sdk-go/types"

func New(ac types.AbstractClient) Tx {
	return txClient{
		AbstractClient: ac,
	}
}

func (t txClient) BuildAndSend(msgs []types.Msg, baseTx types.BaseTx) (types.Result, error) {
	return t.AbstractClient.Broadcast(baseTx, msgs)
}

func (t txClient) Broadcast(signedTx types.StdTx, mode types.BroadcastMode) (types.Result, error) {
	return t.AbstractClient.BroadcastTx(signedTx, mode)
}

func (t txClient) Sign(stdTx types.StdTx, name string, password string, online bool) (types.StdTx, error) {
	return t.AbstractClient.Sign(stdTx, name, password, online)
}
