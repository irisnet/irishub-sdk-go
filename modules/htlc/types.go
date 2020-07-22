package htlc

import (
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

type Htlc struct {
	sender               sdk.AccAddress
	to                   sdk.AccAddress
	receiverOnOtherChain string
	amount               sdk.Coins
	secret               bytes.HexBytes
	timestamp            uint64
	expirationHeight     uint64
	state                int32
}
