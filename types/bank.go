package types

// expose bank module api for user
type Bank interface {
	QueryAccount(address string) (BaseAccount, error)
	QueryTokenStats(tokenID string) (TokenStats, error)
	Send(to string, amount Coins, baseTx BaseTx) (Result, error)
	Burn(amount Coins, baseTx BaseTx) (Result, error)
	SetMemoRegexp(memoRegexp string, baseTx BaseTx) (Result, error)
}

type TokenStats struct {
	LooseTokens  Coins `json:"loose_tokens"`
	BondedTokens Coins `json:"bonded_tokens"`
	BurnedTokens Coins `json:"burned_tokens"`
	TotalSupply  Coins `json:"total_supply"`
}

type EventDataMsgSend struct {
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount []Coin `json:"amount"`
}
type EventMsgSendCallback func(EventDataMsgSend)
