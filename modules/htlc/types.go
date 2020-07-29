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
	sender               sdk.AccAddress
	to                   sdk.AccAddress
	receiverOnOtherChain string
	amount               sdk.Coins
	secret               bytes.HexBytes
	timestamp            uint64
	expirationHeight     uint64
	state                int32
}

func (h htlc) Convert() interface{} {
	return rpc.HTLC{
		Sender:               h.sender.String(),
		To:                   h.to.String(),
		ReceiverOnOtherChain: h.receiverOnOtherChain,
		Secret:               h.secret,
		Timestamp:            h.timestamp,
		ExpirationHeight:     h.expirationHeight,
		State:                h.state,
	}
}
