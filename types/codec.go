package types

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	//"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/sr25519"
)

var codec Codec

func init() {
	codec = NewAminoCodec()
}

type Codec interface {
	MarshalJSON(o interface{}) ([]byte, error)
	UnmarshalJSON(bz []byte, ptr interface{}) error

	MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error)
	UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error

	RegisterConcrete(o interface{}, name string)
	RegisterInterface(ptr interface{})

	MustMarshalBinaryBare(o interface{}) []byte
	MustUnmarshalBinaryBare(bz []byte, ptr interface{})

	UnmarshalBinaryBare(bz []byte, ptr interface{}) error
}

type AminoCodec struct {
	*amino.Codec
}

func NewAminoCodec() Codec {
	cdc := amino.NewCodec()
	return AminoCodec{cdc}
}

func (cdc AminoCodec) RegisterConcrete(o interface{}, name string) {
	cdc.Codec.RegisterConcrete(o, name, nil)
}

func (cdc AminoCodec) RegisterInterface(ptr interface{}) {
	cdc.Codec.RegisterInterface(ptr, nil)
}

func RegisterCodec(cdc Codec) {
	cdc.RegisterInterface((*AccountI)(nil))
	cdc.RegisterInterface((*Msg)(nil))
	//cdc.RegisterConcrete(&BaseAccount{}, "irishub/bank/Account")
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/BaseAccount")
	cdc.RegisterConcrete(StdTx{}, "irishub/bank/StdTx")
	// These are all written here instead of
	cdc.RegisterInterface((*crypto.PubKey)(nil))
	cdc.RegisterConcrete(ed25519.PubKey{},
		ed25519.PubKeyName)
	cdc.RegisterConcrete(sr25519.PubKey{},
		sr25519.PubKeyName)
	cdc.RegisterConcrete(secp256k1.PubKey{},
		secp256k1.PubKeyName)
	//cdc.RegisterConcrete(multisig.PubKeyMultisigThreshold{},
	//	multisig.PubKeyMultisigThresholdAminoRoute)

	cdc.RegisterInterface((*crypto.PrivKey)(nil))
	cdc.RegisterConcrete(ed25519.PrivKey{},
		ed25519.PubKeyName)
	cdc.RegisterConcrete(secp256k1.PrivKey{},
		secp256k1.PubKeyName)
	cdc.RegisterConcrete(sr25519.PrivKey{},
		sr25519.PubKeyName)

	cdc.RegisterInterface((*Store)(nil))
	cdc.RegisterConcrete(PrivKeyInfo{}, "sdk/PrivKeyInfo")
	cdc.RegisterConcrete(KeystoreInfo{}, "sdk/KeystoreInfo")

	codec = cdc
}
