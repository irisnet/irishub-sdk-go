package htlc

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/tendermint/tendermint/libs/bytes"
)

const (
	ModuleName = "htlc"
)

var (
	cdc = original.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

func registerCodec(cdc original.Codec) {

}

type htlc struct {
	Sender           original.AccAddress `json:"sender"`
	To               original.AccAddress `json:"to"`
	Amount           original.Coins      `json:"amount"`
	Secret           bytes.HexBytes      `json:"secret"`
	ExpirationHeight uint64              `json:"expiration_height"`
	State            int32               `json:"state"`
}

func (h htlc) Convert() interface{} {
	return rpc.HTLC{
		Sender:           h.Sender.String(),
		To:               h.To.String(),
		Amount:           h.Amount,
		Secret:           h.Secret.String(),
		ExpirationHeight: h.ExpirationHeight,
		State:            h.State,
	}
}
