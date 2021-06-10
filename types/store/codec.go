package store

import (
	codec2 "github.com/irisnet/irishub-sdk-go/types/codec"
	"github.com/tendermint/tendermint/crypto"

	cryptoAmino "github.com/irisnet/irishub-sdk-go/crypto/codec"
	"github.com/irisnet/irishub-sdk-go/crypto/hd"
)

var cdc *codec2.LegacyAmino

func init() {
	cdc = codec2.NewLegacyAmino()
	cryptoAmino.RegisterCrypto(cdc)
	RegisterCodec(cdc)
	cdc.Seal()
}

// RegisterCodec registers concrete types and interfaces on the given codec.
func RegisterCodec(cdc *codec2.LegacyAmino) {
	cdc.RegisterInterface((*Info)(nil), nil)
	cdc.RegisterConcrete(hd.BIP44Params{}, "crypto/keys/hd/BIP44Params", nil)
	cdc.RegisterConcrete(localInfo{}, "crypto/keys/localInfo", nil)
}

// PubKeyFromBytes unmarshals public key bytes and returns a PubKey
func PubKeyFromBytes(pubKeyBytes []byte) (pubKey crypto.PubKey, err error) {
	err = cdc.UnmarshalBinaryBare(pubKeyBytes, &pubKey)
	return
}
