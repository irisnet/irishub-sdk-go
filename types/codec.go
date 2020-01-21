package utils

import (
	amino "github.com/tendermint/go-amino"
)

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
	return Amino{amino.NewCodec()}
}
