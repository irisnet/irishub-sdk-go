package htlc

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/libs/bytes"
)

const (
	ModuleName = "htlc"
)

var (
	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

func registerCodec(cdc sdk.Codec) {

}

type htlc struct {
	Sender           sdk.AccAddress `json:"sender"`
	To               sdk.AccAddress `json:"to"`
	Amount           sdk.Coins      `json:"amount"`
	Secret           bytes.HexBytes `json:"secret"`
	ExpirationHeight uint64         `json:"expiration_height"`
	State            int32          `json:"state"`
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
