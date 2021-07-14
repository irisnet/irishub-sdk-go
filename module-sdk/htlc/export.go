package htlc

import sdk "github.com/irisnet/core-sdk-go/types"

// expose HTLC module api for user
type Client interface {
	sdk.Module

	CreateHTLC(request CreateHTLCRequest, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)
	ClaimHTLC(hashLock string, secret string, baseTx sdk.BaseTx) (sdk.ResultTx, sdk.Error)

	QueryHTLC(hashLock string) (QueryHTLCResp, sdk.Error)
	QueryParams() (QueryParamsResp, sdk.Error)
}

type CreateHTLCRequest struct {
	To                   string       `json:"to"`
	ReceiverOnOtherChain string       `json:"receiver_on_other_chain"`
	SenderOnOtherChain   string       `json:"sender_on_other_chain"`
	Amount               sdk.DecCoins `json:"amount"`
	HashLock             string       `json:"hash_lock"`
	Timestamp            uint64       `json:"timestamp"`
	TimeLock             uint64       `json:"time_lock"`
	Transfer             bool         ` json:"transfer"`
}

type QueryHTLCResp struct {
	Sender               string    `json:"sender"`
	To                   string    `json:"to"`
	ReceiverOnOtherChain string    `json:"receiver_on_other_chain"`
	SenderOnOtherChain   string    `json:"sender_on_other_chain"`
	Amount               sdk.Coins `json:"amount"`
	Secret               string    `json:"secret"`
	Timestamp            uint64    `json:"timestamp"`
	ExpirationHeight     uint64    `json:"expiration_height"`
	State                int32     `json:"state"`
	Transfer             bool      ` json:"transfer"`
}

type QueryParamsResp struct {
	AssetParams []AssetParamDto `json:"asset_params"`
}

type AssetParamDto struct {
	Denom         string         `json:"denom"`
	SupplyLimit   SupplyLimitDto `json:"supply_limit"`
	Active        bool           `json:"active"`
	DeputyAddress string         `json:"deputy_address"`
	FixedFee      uint64         `json:"fixed_fee"`
	MinSwapAmount uint64         `json:"min_swap_amount"`
	MaxSwapAmount uint64         `json:"max_swap_amount"`
	MinBlockLock  uint64         `json:"min_block_lock"`
	MaxBlockLock  uint64         `json:"max_block_lock"`
}

type SupplyLimitDto struct {
	Limit          uint64 `json:"limit"`
	TimeLimited    bool   `json:"time_limited"`
	TimePeriod     int64  `json:"time_period"`
	TimeBasedLimit uint64 `json:"time_based_limit"`
}
