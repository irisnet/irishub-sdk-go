package coinswap

import (
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/query"
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

	BuyTokenWithAutoEstimate(paidTokenDenom string, boughtCoin sdk.Coin, pools []sdk.PoolInfo,
		deadline int64,
		baseTx sdk.BaseTx,
	) (res *SwapCoinResponse, err error)
	SellTokenWithAutoEstimate(gotTokenDenom string, soldCoin sdk.Coin, pools []sdk.PoolInfo,
		deadline int64,
		baseTx sdk.BaseTx,
	) (res *SwapCoinResponse, err error)

	QueryPool(lptDenom string) (*QueryPoolResponse, error)
	QueryAllPools(pageReq sdk.PageRequest) (*QueryAllPoolsResponse, error)

	EstimateTokenForSoldBase(tokenDenom string,
		soldBase sdk.Int, pool *sdk.PoolInfo,
	) (sdk.Int, error)
	EstimateBaseForSoldToken(soldToken sdk.Coin, pool *sdk.PoolInfo) (sdk.Int, error)
	EstimateTokenForSoldToken(boughtTokenDenom string,
		soldToken sdk.Coin, pools []sdk.PoolInfo) (sdk.Int, error)
	EstimateTokenForBoughtBase(soldTokenDenom string,
		boughtBase sdk.Int, pool *sdk.PoolInfo) (sdk.Int, error)
	EstimateBaseForBoughtToken(boughtToken sdk.Coin, pool *sdk.PoolInfo) (sdk.Int, error)
	EstimateTokenForBoughtToken(soldTokenDenom string,
		boughtToken sdk.Coin, pool []sdk.PoolInfo) (sdk.Int, error)
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
	Liquidity sdk.Coin
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
	Pool sdk.PoolInfo
}

type QueryAllPoolsResponse struct {
	Pools      []sdk.PoolInfo
	Pagination *query.PageResponse
}
