package random

import sdk "github.com/irisnet/core-sdk-go/types"

// expose Random module api for user
type Client interface {
	sdk.Module

	RequestRandom(request RequestRandomRequest, basTx sdk.BaseTx) (RequestRandomResp, sdk.ResultTx, sdk.Error)

	QueryRandom(ReqId string) (QueryRandomResp, sdk.Error)
	QueryRandomRequestQueue(height int64) ([]QueryRandomRequestQueueResp, sdk.Error)
}

type RequestRandomRequest struct {
	BlockInterval uint64    `json:"block_interval"`
	Oracle        bool      `json:"oracle"`
	ServiceFeeCap sdk.Coins `json:"service_fee_cap"`
}

type RequestRandomResp struct {
	Height int64  `json:"height"`
	ReqID  string `json:"req_id"`
}

type QueryRandomResp struct {
	RequestTxHash string `json:"request_tx_hash" yaml:"request_tx_hash"`
	Height        int64  `json:"height" yaml:"height"`
	Value         string `json:"value" yaml:"value"`
}

type QueryRandomRequestQueueResp struct {
	Height           int64     `json:"height" yaml:"height"`
	Consumer         string    `json:"consumer" yaml:"consumer"`
	TxHash           string    `json:"tx_hash" yaml:"tx_hash"`
	Oracle           bool      `json:"oracle" yaml:"oracle"`
	ServiceFeeCap    sdk.Coins `json:"service_fee_cap" yaml:"service_fee_cap"`
	ServiceContextId string    `json:"service_context_id" yaml:"service_context_id"`
}
