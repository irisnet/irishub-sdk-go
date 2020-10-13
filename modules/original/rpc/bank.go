package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

// expose bank module api for user
type Bank interface {
	original.Module
	QueryBalances(address, denom string) (original.Balances, original.Error)
	QueryAccount(address string) (original.BaseAccount, original.Error)
	QueryTokenStats(tokenID string) (TokenStats, original.Error)
	QueryTotalSupply() (original.Coins, original.Error)
	Send(to string, amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error)
	MultiSend(receipts Receipts, baseTx original.BaseTx) ([]original.ResultTx, original.Error)
	Burn(amount original.DecCoins, baseTx original.BaseTx) (original.ResultTx, original.Error)
	SetMemoRegexp(memoRegexp string, baseTx original.BaseTx) (original.ResultTx, original.Error)
	SubscribeSendTx(from, to string, callback EventMsgSendCallback) original.Subscription
}

type Receipt struct {
	Address string            `json:"address"`
	Amount  original.DecCoins `json:"amount"`
}

type Receipts []Receipt

func (r Receipts) Len() int {
	return len(r)
}

func (r Receipts) Sub(begin, end int) original.SplitAble {
	return r[begin:end]
}

type TokenStats struct {
	LooseTokens  original.Coins `json:"loose_tokens"`
	BondedTokens original.Coins `json:"bonded_tokens"`
	BurnedTokens original.Coins `json:"burned_tokens"`
	TotalSupply  original.Coins `json:"total_supply"`
}

type EventDataMsgSend struct {
	Height int64           `json:"height"`
	Hash   string          `json:"hash"`
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount []original.Coin `json:"amount"`
}
type EventMsgSendCallback func(EventDataMsgSend)
