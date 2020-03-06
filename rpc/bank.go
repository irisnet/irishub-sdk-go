package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

// expose bank module api for user
type Bank interface {
	sdk.Module
	QueryAccount(address string) (sdk.BaseAccount, error)
	QueryTokenStats(tokenID string) (TokenStats, error)
	Send(to string, amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error)
	Burn(amount sdk.Coins, baseTx sdk.BaseTx) (sdk.Result, error)
	SetMemoRegexp(memoRegexp string, baseTx sdk.BaseTx) (sdk.Result, error)
	SubscribeSendTx(from, to string, callback EventMsgSendCallback) sdk.Subscription
}

type TokenStats struct {
	LooseTokens  sdk.Coins `json:"loose_tokens"`
	BondedTokens sdk.Coins `json:"bonded_tokens"`
	BurnedTokens sdk.Coins `json:"burned_tokens"`
	TotalSupply  sdk.Coins `json:"total_supply"`
}

type EventDataMsgSend struct {
	Height int64      `json:"height"`
	Hash   string     `json:"hash"`
	From   string     `json:"from"`
	To     string     `json:"to"`
	Amount []sdk.Coin `json:"amount"`
}
type EventMsgSendCallback func(EventDataMsgSend)
