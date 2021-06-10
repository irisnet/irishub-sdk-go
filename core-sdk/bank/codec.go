package bank

import (
	"github.com/irisnet/irishub-sdk-go/auth"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

var (
	amino     = sdk.NewLegacyAmino()
	ModuleCdc = sdk.NewAminoCodec(amino)
)

func init() {
	sdk.RegisterCrypto(amino)
	amino.Seal()
}

func RegisterInterfaces(registry sdk.InterfaceRegistry) {
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
