package rpc

import "github.com/irisnet/irishub-sdk-go/types"

// expose bank module api for user
type Bank interface {
	types.Module
	QueryAccount(address string) (types.BaseAccount, error)
	QueryTokenStats(tokenID string) (TokenStats, error)
	Send(to string, amount types.Coins, baseTx types.BaseTx) (types.Result, error)
	Burn(amount types.Coins, baseTx types.BaseTx) (types.Result, error)
	SetMemoRegexp(memoRegexp string, baseTx types.BaseTx) (types.Result, error)
	SubscribeSendTx(from, to string, callback EventMsgSendCallback) types.Subscription
}

type TokenStats struct {
	LooseTokens  types.Coins `json:"loose_tokens"`
	BondedTokens types.Coins `json:"bonded_tokens"`
	BurnedTokens types.Coins `json:"burned_tokens"`
	TotalSupply  types.Coins `json:"total_supply"`
}

type EventDataMsgSend struct {
	Height int64        `json:"height"`
	Hash   string       `json:"hash"`
	From   string       `json:"from"`
	To     string       `json:"to"`
	Amount []types.Coin `json:"amount"`
}
type EventMsgSendCallback func(EventDataMsgSend)
