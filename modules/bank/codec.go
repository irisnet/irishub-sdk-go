package bank

import (
	"github.com/irisnet/irishub-sdk-go/codec"
	cryptocodec "github.com/irisnet/irishub-sdk-go/crypto/codec"
	"github.com/irisnet/irishub-sdk-go/modules/auth"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	codec2 "github.com/irisnet/irishub-sdk-go/types/codec"
	"github.com/irisnet/irishub-sdk-go/types/codec/types"
)

var (
	amino     = codec2.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSend{},
		&MsgMultiSend{},
	)

	registry.RegisterImplementations(
		(*auth.Account)(nil),
		&auth.BaseAccount{},
	)
}
