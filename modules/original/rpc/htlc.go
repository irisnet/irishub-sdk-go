package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

type Htlc interface {
	original.Module
	QueryHtlc(hashLock string) (HTLC, original.Error)
}

type HTLC struct {
	Sender           string
	To               string
	Amount           original.Coins
	Secret           string
	ExpirationHeight uint64
	State            int32
}
