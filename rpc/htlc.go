package rpc

import sdk "github.com/irisnet/irishub-sdk-go/types"

type Htlc interface {
	sdk.Module
	QueryHTLC(hashLock string) (interface{}, sdk.Error)
}
