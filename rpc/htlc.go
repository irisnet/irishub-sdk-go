package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/tendermint/tendermint/libs/bytes"
)

type Htlc interface {
	sdk.Module
	QueryHtlc(hashLock string) (HTLC, sdk.Error)
}

type HTLC struct {
	Sender               string
	To                   string
	ReceiverOnOtherChain string
	Amount               sdk.Coins
	Secret               bytes.HexBytes
	Timestamp            uint64
	ExpirationHeight     uint64
	State                int32
}
