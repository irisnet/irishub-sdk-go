package bank

import (
	commoncodec "github.com/irisnet/irishub-sdk-go/common/codec"
	"github.com/irisnet/irishub-sdk-go/common/codec/types"
	commoncryptocodec "github.com/irisnet/irishub-sdk-go/common/crypto/codec"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/auth"
)

var (
	amino     = commoncodec.NewLegacyAmino()
	ModuleCdc = commoncodec.NewAminoCodec(amino)
)

func init() {
	commoncryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

// No duplicate registration
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
