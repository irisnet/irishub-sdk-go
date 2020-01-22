package types

import (
	"github.com/tendermint/go-amino"
)

var defaultCdc Codec

func init() {
	defaultCdc = NewAmino()
}

type Encode interface {
	MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error)
	MarshalJSON(o interface{}) ([]byte, error)
}

type Decode interface {
	UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error
	UnmarshalJSON(bz []byte, ptr interface{}) error
}

type Codec interface {
	Encode
	Decode
}

type Amino struct {
	*amino.Codec
}

func NewAmino() Amino {
	cdc := amino.NewCodec()
	RegisterAuth(cdc)
	RegisterBank(cdc)
	RegisterStake(cdc)
	return Amino{cdc}
}

func DefaultCodec() Codec {
	return defaultCdc
}
