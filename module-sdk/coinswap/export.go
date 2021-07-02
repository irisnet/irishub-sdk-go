package coinswap

import (
	sdk "github.com/irisnet/core-sdk-go/types"
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

	BuyTokenWithAutoEstimate(paidTokenDenom string, boughtCoin sdk.Coin,
		deadline int64,
		baseTx sdk.BaseTx,
	) (res *SwapCoinResponse, err error)
	SellTokenWithAutoEstimate(gotTokenDenom string, soldCoin sdk.Coin,
		deadline int64,
		baseTx sdk.BaseTx,
	) (res *SwapCoinResponse, err error)

	QueryPool(denom string) (*QueryPoolResponse, error)
	QueryAllPools() (*QueryAllPoolsResponse, error)

	EstimateTokenForSoldBase(tokenDenom string,
		soldBase sdk.Int,
	) (sdk.Int, error)
	EstimateBaseForSoldToken(soldToken sdk.Coin) (sdk.Int, error)
	EstimateTokenForSoldToken(boughtTokenDenom string,
		soldToken sdk.Coin) (sdk.Int, error)
	EstimateTokenForBoughtBase(soldTokenDenom string,
		boughtBase sdk.Int) (sdk.Int, error)
	EstimateBaseForBoughtToken(boughtToken sdk.Coin) (sdk.Int, error)
	EstimateTokenForBoughtToken(soldTokenDenom string,
		boughtToken sdk.Coin) (sdk.Int, error)
}

type AddLiquidityRequest struct {
	MaxToken     sdk.Coin
	BaseAmt      sdk.Int
	MinLiquidity sdk.Int
	Deadline     int64
}

type AddLiquidityResponse struct {
	TokenAmt  sdk.Int
	BaseAmt   sdk.Int
	Liquidity sdk.Int
	TxHash    string
}

type RemoveLiquidityRequest struct {
	MinTokenAmt sdk.Int
	MinBaseAmt  sdk.Int
	Liquidity   sdk.Coin
	Deadline    int64
}

type RemoveLiquidityResponse struct {
	TokenAmt  sdk.Int
	BaseAmt   sdk.Int
	Liquidity sdk.Coin
	TxHash    string
}

type SwapCoinRequest struct {
	Input      sdk.Coin
	Output     sdk.Coin
	Receiver   string
	Deadline   int64
	IsBuyOrder bool
}

type SwapCoinResponse struct {
	InputAmt  sdk.Int
	OutputAmt sdk.Int
	TxHash    string
}

type QueryPoolResponse struct {
	BaseCoin  sdk.Coin
	TokenCoin sdk.Coin
	Liquidity sdk.Coin
	Fee       string
}

type QueryAllPoolsResponse struct {
	Pools []QueryPoolResponse
}
