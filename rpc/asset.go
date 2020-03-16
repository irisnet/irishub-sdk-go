package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Asset interface {
	sdk.Module
	QueryToken(symbol string) (sdk.Token, error)
}
