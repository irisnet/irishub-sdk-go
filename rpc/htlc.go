package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Htlc interface {
	sdk.Module
	QueryHtlc(hashLock string) (HTLC, sdk.Error)
}

type HTLC struct {
	Sender           string
	To               string
	Amount           sdk.Coins
	Secret           string
	ExpirationHeight uint64
	State            int32
}
