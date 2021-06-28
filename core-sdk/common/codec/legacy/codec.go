package legacy

import (
	commoncodec "github.com/irisnet/irishub-sdk-go/common/codec"

	cryptocodec "github.com/irisnet/irishub-sdk-go/common/crypto/codec"
)

// Cdc defines a global generic sealed Amino codec to be used throughout sdk. It
// has all Tendermint crypto and evidence types registered.
//
// TODO: Deprecated - remove this global.
var Cdc *commoncodec.LegacyAmino

func init() {
	Cdc = commoncodec.NewLegacyAmino()
	cryptocodec.RegisterCrypto(Cdc)
	commoncodec.RegisterEvidences(Cdc)
	Cdc.Seal()
}
