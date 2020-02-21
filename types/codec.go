package types

import (
	"github.com/tendermint/go-amino"
)

type Codec interface {
	MarshalJSON(o interface{}) ([]byte, error)
	UnmarshalJSON(bz []byte, ptr interface{}) error

	MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error)
	UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error

	RegisterConcrete(o interface{}, name string)
	RegisterInterface(ptr interface{})
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
