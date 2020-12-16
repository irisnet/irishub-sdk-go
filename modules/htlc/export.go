package htlc

import sdk "github.com/irisnet/irishub-sdk-go/types"

// expose HTLC module api for user
type Client interface {
	sdk.Module

	CreateHTLC(request CreateHTLCRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	ClaimHTLC(hashLock string, secret string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	RefundHTLC(hashLock string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryHTLC(hashLock string) (QueryHTLCResp, sdk.Error)
}

type CreateHTLCRequest struct {
	To                   string       `json:"to"`
	ReceiverOnOtherChain string       `json:"receiver_on_other_chain"`
	Amount               sdk.DecCoins `json:"amount"`
	HashLock             string       `json:"hash_lock"`
	Timestamp            uint64       `json:"timestamp"`
	TimeLock             uint64       `json:"time_lock"`
}

type QueryHTLCResp struct {
	Sender               string    `json:"sender"`
	To                   string    `json:"to"`
	ReceiverOnOtherChain string    `json:"receiver_on_other_chain"`
	Amount               sdk.Coins `json:"amount"`
	Secret               string    `json:"secret"`
	Timestamp            uint64    `json:"timestamp"`
	ExpirationHeight     uint64    `json:"expiration_height"`
	State                int32     `json:"state"`
}
