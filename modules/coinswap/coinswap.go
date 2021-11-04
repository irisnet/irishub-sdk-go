package coinswap

import (
	"context"
	"errors"
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types/query"

	ctypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub-sdk-go/codec"
	"github.com/irisnet/irishub-sdk-go/codec/types"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type coinswapClient struct {
	sdk.BaseClient
	codec.Marshaler
	totalSupply
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler, queryTotalSupply totalSupply) Client {
	return coinswapClient{
		BaseClient:  bc,
		Marshaler:   cdc,
		totalSupply: queryTotalSupply,
	}
}

func (swap coinswapClient) Name() string {
	return ModuleName
}

func (swap coinswapClient) RegisterInterfaceTypes(registry types.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (swap coinswapClient) AddLiquidity(request AddLiquidityRequest,
	baseTx sdk.BaseTx) (*AddLiquidityResponse, error) {
	creator, err := swap.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	msg := &MsgAddLiquidity{
		MaxToken:         request.MaxToken,
		ExactStandardAmt: ctypes.NewInt(request.BaseAmt.Int64()),
		MinLiquidity:     ctypes.NewInt(request.MinLiquidity.Int64()),
		Deadline:         request.Deadline,
		Sender:           creator.String(),
	}

	res, err := swap.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return nil, err
	}

	var totalCoins = sdk.NewCoins()
	coinStrs := res.Events.GetValues(eventTypeTransfer, attributeKeyAmount)
	for _, coinStr := range coinStrs {
		coins, er := sdk.ParseCoins(coinStr)
		if er != nil {
			swap.Logger().Error("Parse coin str failed", "coin", coinStr)
			continue
		}
		totalCoins = totalCoins.Add(coins...)
	}
	lptCoin, er := sdk.ParseCoin(res.Events.GetValues(eventTypeCoinbase, attributeKeyAmount)[0])
	if er != nil {
		return nil, er
	}
	response := &AddLiquidityResponse{
		TokenAmt:  totalCoins.AmountOf(request.MaxToken.Denom),
		BaseAmt:   request.BaseAmt,
		Liquidity: lptCoin,
		TxHash:    res.Hash,
	}
	return response, nil
}

func (swap coinswapClient) RemoveLiquidity(request RemoveLiquidityRequest,
	baseTx sdk.BaseTx) (*RemoveLiquidityResponse, error) {
	creator, err := swap.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	msg := &MsgRemoveLiquidity{
		WithdrawLiquidity: request.Liquidity,
		MinToken:          ctypes.NewInt(request.MinTokenAmt.Int64()),
		MinStandardAmt:    ctypes.NewInt(request.MinBaseAmt.Int64()),
		Deadline:          request.Deadline,
		Sender:            creator.String(),
	}

	res, err := swap.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return nil, err
	}

	var totalCoins = sdk.NewCoins()
	coinStrs := res.Events.GetValues(eventTypeTransfer, attributeKeyAmount)
	for _, coinStr := range coinStrs {
		coins, er := sdk.ParseCoins(coinStr)
		if er != nil {
			swap.Logger().Error("Parse coin str failed", "coin", coinStr)
			continue
		}
		totalCoins = totalCoins.Add(coins...)
	}
	pool, er := swap.QueryPool(request.Liquidity.Denom)
	if er != nil {
		return nil, er
	}
	response := &RemoveLiquidityResponse{
		TokenAmt:  totalCoins.AmountOf(pool.Pool.Token.Denom),
		BaseAmt:   totalCoins.AmountOf(sdk.BaseDenom),
		Liquidity: request.Liquidity,
		TxHash:    res.Hash,
	}
	return response, nil
}

func (swap coinswapClient) SwapCoin(request SwapCoinRequest, baseTx sdk.BaseTx) (*SwapCoinResponse, error) {
	creator, err := swap.QueryAddress(baseTx.From, baseTx.Password)
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	input := Input{
		Address: creator.String(),
		Coin:    request.Input,
	}

	if len(request.Receiver) == 0 {
		request.Receiver = input.Address
	}

	output := Output{
		Address: request.Receiver,
		Coin:    request.Output,
	}

	msg := &MsgSwapOrder{
		Input:      input,
		Output:     output,
		Deadline:   request.Deadline,
		IsBuyOrder: request.IsBuyOrder,
	}

	res, err := swap.BuildAndSend([]sdk.Msg{msg}, baseTx)
	if err != nil {
		return nil, err
	}

	amount, er := res.Events.GetValue(eventTypeSwap, attributeKeyAmount)
	if er != nil {
		return nil, er
	}

	amt, ok := sdk.NewIntFromString(amount)
	if !ok {
		return nil, sdk.Wrapf("%s can not convert to sdk.Int type", amount)
	}

	inputAmt := request.Input.Amount
	outputAmt := request.Output.Amount
	if request.IsBuyOrder {
		inputAmt = amt
	} else {
		outputAmt = amt
	}

	response := &SwapCoinResponse{
		InputAmt:  inputAmt,
		OutputAmt: outputAmt,
		TxHash:    res.Hash,
	}
	return response, nil
}

func (swap coinswapClient) BuyTokenWithAutoEstimate(paidTokenDenom string, 
	boughtCoin sdk.Coin, pools []sdk.PoolInfo, deadline int64, baseTx sdk.BaseTx,
) (res *SwapCoinResponse, err error) {
	var pool *sdk.PoolInfo
	var amount = sdk.ZeroInt()
	switch {
	case paidTokenDenom == sdk.BaseDenom:
		pool, err = FilterPoolByDenom(pools, boughtCoin.Denom)
		if err == nil {
			amount, err = swap.EstimateBaseForBoughtToken(boughtCoin, pool)
		}
		break
	case boughtCoin.Denom == sdk.BaseDenom:
		pool, err = FilterPoolByDenom(pools, paidTokenDenom)
		if err == nil {
			amount, err = swap.EstimateTokenForBoughtBase(paidTokenDenom, boughtCoin.Amount, pool)
		}
		break
	default:
		amount, err = swap.EstimateTokenForBoughtToken(paidTokenDenom, boughtCoin, pools)
		break
	}
	if err != nil {
		return nil, err
	}
	req := SwapCoinRequest{
		Input:      sdk.NewCoin(paidTokenDenom, amount),
		Output:     boughtCoin,
		Deadline:   deadline,
		IsBuyOrder: true,
	}
	return swap.SwapCoin(req, baseTx)
}

func (swap coinswapClient) SellTokenWithAutoEstimate(gotTokenDenom string, 
	soldCoin sdk.Coin, pools []sdk.PoolInfo, deadline int64, baseTx sdk.BaseTx,
) (res *SwapCoinResponse, err error) {
	var pool *sdk.PoolInfo
	var amount = sdk.ZeroInt()
	switch {
	case gotTokenDenom == sdk.BaseDenom:
		pool, err = FilterPoolByDenom(pools, soldCoin.Denom)
		if err == nil {
			amount, err = swap.EstimateBaseForSoldToken(soldCoin, pool)
		}
		break
	case soldCoin.Denom == sdk.BaseDenom:
		pool, err = FilterPoolByDenom(pools, gotTokenDenom)
		if err == nil {
			amount, err = swap.EstimateTokenForSoldBase(gotTokenDenom, soldCoin.Amount, pool)
		}
		break
	default:
		amount, err = swap.EstimateTokenForSoldToken(gotTokenDenom, soldCoin, pools)
		break
	}

	if err != nil {
		return nil, err
	}

	req := SwapCoinRequest{
		Input:      soldCoin,
		Output:     sdk.NewCoin(gotTokenDenom, amount),
		Deadline:   deadline,
		IsBuyOrder: false,
	}
	return swap.SwapCoin(req, baseTx)
}

func (swap coinswapClient) EstimateTokenForSoldBase(tokenDenom string,
	soldBaseAmt sdk.Int, pool *sdk.PoolInfo) (sdk.Int, error) {
	if tokenDenom != pool.Token.Denom {
		return sdk.ZeroInt(), errors.New(
			fmt.Sprintf("bought token %s and pool token %s should be same", 
				tokenDenom, pool.Token.Denom))
	}
	fee := sdk.MustNewDecFromStr(pool.Fee)
	amount := getInputPrice(soldBaseAmt,
		pool.Standard.Amount, pool.Token.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateBaseForSoldToken(soldToken sdk.Coin,
	pool *sdk.PoolInfo) (sdk.Int, error) {
	if soldToken.Denom != pool.Token.Denom {
		return sdk.ZeroInt(), errors.New(
			fmt.Sprintf("bought token %s and pool token %s should be same", 
				soldToken.Denom, pool.Token.Denom))
	}
	fee := sdk.MustNewDecFromStr(pool.Fee)
	amount := getInputPrice(soldToken.Amount,
		pool.Token.Amount, pool.Standard.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateTokenForSoldToken(boughtTokenDenom string,
	soldToken sdk.Coin, pools []sdk.PoolInfo) (sdk.Int, error) {
	if boughtTokenDenom == soldToken.Denom {
		return sdk.ZeroInt(), errors.New("invalid trade")
	}
	pool, err := FilterPoolByDenom(pools, soldToken.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	boughtBaseAmt, err := swap.EstimateBaseForSoldToken(soldToken, pool)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	pool, err = FilterPoolByDenom(pools, boughtTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return swap.EstimateTokenForSoldBase(boughtTokenDenom, boughtBaseAmt, pool)
}

func (swap coinswapClient) EstimateTokenForBoughtBase(soldTokenDenom string,
	exactBoughtBaseAmt sdk.Int, pool *sdk.PoolInfo) (sdk.Int, error) {
	if soldTokenDenom != pool.Token.Denom {
		return sdk.ZeroInt(), errors.New(
			fmt.Sprintf("sold token %s and pool token %s should be same", 
				soldTokenDenom, pool.Token.Denom))
	}
	fee := sdk.MustNewDecFromStr(pool.Fee)
	amount := getOutputPrice(exactBoughtBaseAmt,
		pool.Token.Amount, pool.Standard.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateBaseForBoughtToken(boughtToken sdk.Coin, 
	pool *sdk.PoolInfo) (sdk.Int, error) {
	if boughtToken.Denom != pool.Token.Denom {
		return sdk.ZeroInt(), errors.New(
			fmt.Sprintf("bought token %s and pool token %s should be same", 
				boughtToken, pool.Token.Denom))
	}
	fee := sdk.MustNewDecFromStr(pool.Fee)
	amount := getOutputPrice(boughtToken.Amount,
		pool.Standard.Amount, pool.Token.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateTokenForBoughtToken(soldTokenDenom string,
	boughtToken sdk.Coin, pools []sdk.PoolInfo) (sdk.Int, error) {
	if soldTokenDenom == boughtToken.Denom {
		return sdk.ZeroInt(), errors.New("invalid trade")
	}
	pool, err := FilterPoolByDenom(pools, boughtToken.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	soldBaseAmt, err := swap.EstimateBaseForBoughtToken(boughtToken, pool)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	pool, err = FilterPoolByDenom(pools, soldTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return swap.EstimateTokenForBoughtBase(soldTokenDenom, soldBaseAmt, pool)
}

func (swap coinswapClient) QueryPool(lptDenom string) (*QueryPoolResponse, error) {
	conn, err := swap.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).LiquidityPool(
		context.Background(),
		&QueryLiquidityPoolRequest{LptDenom: lptDenom},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return resp.Convert().(*QueryPoolResponse), err
}

func (swap coinswapClient) QueryAllPools(req sdk.PageRequest) (*QueryAllPoolsResponse, error) {
	conn, err := swap.GenConn()
	defer func() { _ = conn.Close() }()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).LiquidityPools(
		context.Background(),
		&QueryLiquidityPoolsRequest{
			Pagination: &query.PageRequest{
				Key:        req.Key,
				Offset:     req.Offset,
				Limit:      req.Limit,
				CountTotal: req.CountTotal,
			},
		},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return resp.Convert().(*QueryAllPoolsResponse), err
}

func FilterPoolByDenom(pools []sdk.PoolInfo, denom string) (*sdk.PoolInfo, error) {
	for _, pool := range pools {
		if pool.Token.Denom == denom {
			return &pool, nil
		}
	}
	return nil, sdk.Wrapf("denom: %s not exist in pools: %v", denom, pools)
}

// getInputPrice returns the amount of coins bought (calculated) given the input amount being sold (exact)
// The fee is included in the input coins being bought
// https://github.com/runtimeverification/verified-smart-contracts/blob/uniswap/uniswap/x-y-k.pdf
func getInputPrice(inputAmt, inputReserve, outputReserve sdk.Int, fee sdk.Dec) sdk.Int {
	deltaFee := sdk.OneDec().Sub(fee)
	inputAmtWithFee := inputAmt.Mul(sdk.NewIntFromBigInt(deltaFee.BigInt()))
	numerator := inputAmtWithFee.Mul(outputReserve)
	denominator := inputReserve.Mul(sdk.NewIntWithDecimal(1, sdk.Precision)).Add(inputAmtWithFee)
	return numerator.Quo(denominator)
}

// getOutputPrice returns the amount of coins sold (calculated) given the output amount being bought (exact)
// The fee is included in the output coins being bought
func getOutputPrice(outputAmt, inputReserve, outputReserve sdk.Int, fee sdk.Dec) sdk.Int {
	deltaFee := sdk.OneDec().Sub(fee)
	numerator := inputReserve.Mul(outputAmt).Mul(sdk.NewIntWithDecimal(1, sdk.Precision))
	denominator := (outputReserve.Sub(outputAmt)).Mul(sdk.NewIntFromBigInt(deltaFee.BigInt()))
	return numerator.Quo(denominator).Add(sdk.OneInt())
}
