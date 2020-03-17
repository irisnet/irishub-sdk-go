package rpc

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type Asset interface {
	sdk.Module
	QueryToken(symbol string) (sdk.Token, error)
	QueryTokens(owner string) (sdk.Tokens, error)
	QueryFees(symbol string) (TokenFees, error)
}

// TokenFees is for the token fees query output
type TokenFees struct {
	Exist    bool     `json:"exist"`     // indicate if the token has existed
	IssueFee sdk.Coin `json:"issue_fee"` // issue fee
	MintFee  sdk.Coin `json:"mint_fee"`  // mint fee
}
