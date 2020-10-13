package rpc

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
)

type Asset interface {
	original.Module
	QueryTokens() (original.Tokens, error)
	QueryTokenDenom(denom string) (original.TokenData, error)
	QueryToken(symbol string) (original.Token, error)
}

// TokenFees is for the token fees query output
type TokenFees struct {
	Exist    bool          `json:"exist"`     // indicate if the token has existed
	IssueFee original.Coin `json:"issue_fee"` // issue fee
	MintFee  original.Coin `json:"mint_fee"`  // mint fee
}
