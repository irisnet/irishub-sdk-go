package asset

import (
	"github.com/irisnet/irishub-sdk-go/rpc"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

const (
	ModuleName = "asset"
)

var (
	cdc = sdk.NewAminoCodec()
)

func init() {
	registerCodec(cdc)
}

func registerCodec(cdc sdk.Codec) {
}

// tokenFees is for the token fees query output
type tokenFees struct {
	Exist    bool     `json:"exist"`     // indicate if the token has existed
	IssueFee sdk.Coin `json:"issue_fee"` // issue fee
	MintFee  sdk.Coin `json:"mint_fee"`  // mint fee
}

func (t tokenFees) Convert() interface{} {
	return rpc.TokenFees{
		Exist:    t.Exist,
		IssueFee: t.IssueFee,
		MintFee:  t.MintFee,
	}
}
