package asset

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "asset"
)

var (
	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

func registerCodec(cdc sdk.Codec) {
}
