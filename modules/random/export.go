package random

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

// expose Random module api for user
type RandomI interface {
	sdk.Module

	QueryRandom(requestID string) (QueryRandomResp, sdk.Error)
}

type QueryRandomResp struct {
	RequestTxHash string `json:"request_tx_hash" yaml:"request_tx_hash"`
	Height        int64  `json:"height" yaml:"height"`
	Value         string `json:"value" yaml:"value"`
}
