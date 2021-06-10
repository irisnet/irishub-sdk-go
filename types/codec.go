package types

import (
	"github.com/irisnet/irishub-sdk-go/codec"
	codec2 "github.com/irisnet/irishub-sdk-go/types/codec"
	"github.com/irisnet/irishub-sdk-go/types/codec/types"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Marshaler
	TxConfig          TxConfig
	Amino             *codec2.LegacyAmino
}
