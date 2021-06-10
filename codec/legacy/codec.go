package legacy

import (
	cryptocodec "github.com/irisnet/irishub-sdk-go/crypto/codec"
	codec2 "github.com/irisnet/irishub-sdk-go/types/codec"
)

// Cdc defines a global generic sealed Amino codec to be used throughout sdk. It
// has all Tendermint crypto and evidence types registered.
//
// TODO: Deprecated - remove this global.
var Cdc *codec2.LegacyAmino

func init() {
	Cdc = codec2.NewLegacyAmino()
	cryptocodec.RegisterCrypto(Cdc)
	codec2.RegisterEvidences(Cdc)
	Cdc.Seal()
}
