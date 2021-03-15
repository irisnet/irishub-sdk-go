package coinswap

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
	types "github.com/irisnet/irishub-sdk-go/types"
)

// expose Record module api for user
type Client interface {
	sdk.Module
	AddLiquidity(request AddLiquidityRequest,
		baseTx sdk.BaseTx) (*AddLiquidityResponse, error)
	RemoveLiquidity(request RemoveLiquidityRequest,
		baseTx sdk.BaseTx) (*RemoveLiquidityResponse, error)
	SwapCoin(request SwapCoinRequest,
		baseTx sdk.BaseTx) (*SwapCoinResponse, error)
	QueryPool(denom string) (*QueryPoolResponse, error)
}

type AddLiquidityRequest struct {
	MaxToken     types.Coin
	BaseAmt      types.Int
	MinLiquidity types.Int
	Deadline     int64
}

type AddLiquidityResponse struct {
	TokenAmt  types.Int
	BaseAmt   types.Int
	Liquidity types.Int
	TxHash    string
}

type RemoveLiquidityRequest struct {
	MinTokenAmt types.Int
	MinBaseAmt  types.Int
	Liquidity   types.Coin
	Deadline    int64
}

type RemoveLiquidityResponse struct {
	TokenAmt  types.Int
	BaseAmt   types.Int
	Liquidity types.Coin
	TxHash    string
}

type SwapCoinRequest struct {
	Input      types.Coin
	Output     types.Coin
	Receiver   string
	Deadline   int64
	IsBuyOrder bool
}

type SwapCoinResponse struct {
	InputAmt  types.Int
	OutputAmt types.Int
	TxHash    string
}

type QueryPoolResponse struct {
	BaseCoin  types.Coin
	TokenCoin types.Coin
	Liquidity types.Coin
}
