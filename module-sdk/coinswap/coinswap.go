package coinswap

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/irisnet/core-sdk-go/common/codec"
	"github.com/irisnet/core-sdk-go/common/codec/types"
	sdk "github.com/irisnet/core-sdk-go/types"
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
		ExactStandardAmt: request.BaseAmt,
		MinLiquidity:     request.MinLiquidity,
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

	liquidityDenom, er := GetLiquidityDenomFrom(request.MaxToken.Denom)
	if er != nil {
		return nil, er
	}
	response := &AddLiquidityResponse{
		TokenAmt:  totalCoins.AmountOf(request.MaxToken.Denom),
		BaseAmt:   request.BaseAmt,
		Liquidity: totalCoins.AmountOf(liquidityDenom),
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
		MinToken:          request.MinTokenAmt,
		MinStandardAmt:    request.MinBaseAmt,
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

	tokenDenom, er := GetTokenDenomFrom(request.Liquidity.Denom)
	if er != nil {
		return nil, er
	}

	response := &RemoveLiquidityResponse{
		TokenAmt:  totalCoins.AmountOf(tokenDenom),
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

func (swap coinswapClient) BuyTokenWithAutoEstimate(paidTokenDenom string, boughtCoin sdk.Coin,
	deadline int64,
	baseTx sdk.BaseTx,
) (res *SwapCoinResponse, err error) {
	var amount = sdk.ZeroInt()
	switch {
	case paidTokenDenom == sdk.BaseDenom:
		amount, err = swap.EstimateBaseForBoughtToken(boughtCoin)
		break
	case boughtCoin.Denom == sdk.BaseDenom:
		amount, err = swap.EstimateTokenForBoughtBase(paidTokenDenom, boughtCoin.Amount)
		break
	default:
		amount, err = swap.EstimateTokenForBoughtToken(paidTokenDenom, boughtCoin)
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

func (swap coinswapClient) SellTokenWithAutoEstimate(gotTokenDenom string, soldCoin sdk.Coin,
	deadline int64,
	baseTx sdk.BaseTx,
) (res *SwapCoinResponse, err error) {
	var amount = sdk.ZeroInt()
	switch {
	case gotTokenDenom == sdk.BaseDenom:
		amount, err = swap.EstimateBaseForSoldToken(soldCoin)
		break
	case soldCoin.Denom == sdk.BaseDenom:
		amount, err = swap.EstimateTokenForSoldBase(gotTokenDenom, soldCoin.Amount)
		break
	default:
		amount, err = swap.EstimateTokenForSoldToken(gotTokenDenom, soldCoin)
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

func (swap coinswapClient) QueryPool(denom string) (*QueryPoolResponse, error) {
	conn, err := swap.GenConn()

	if err != nil {
		return nil, sdk.Wrap(err)
	}

	resp, err := NewQueryClient(conn).Liquidity(
		context.Background(),
		&QueryLiquidityRequest{Denom: denom},
	)
	if err != nil {
		return nil, sdk.Wrap(err)
	}
	return resp.Convert().(*QueryPoolResponse), err
}

func (swap coinswapClient) QueryAllPools() (*QueryAllPoolsResponse, error) {
	coins, err := swap.totalSupply()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	var pools []QueryPoolResponse
	for _, coin := range coins {
		//Compatible with old data
		if strings.HasPrefix(coin.Denom, "swap/") {
			continue
		}
		denom, err := GetTokenDenomFrom(coin.Denom)
		if err != nil {
			continue
		}
		res, err := swap.QueryPool(denom)
		if err != nil {
			return nil, sdk.Wrap(err)
		}
		pools = append(pools, *res)
	}
	return &QueryAllPoolsResponse{pools}, err
}

func (swap coinswapClient) EstimateTokenForSoldBase(tokenDenom string,
	soldBaseAmt sdk.Int,
) (sdk.Int, error) {
	result, err := swap.QueryPool(tokenDenom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	fee := sdk.MustNewDecFromStr(result.Fee)
	amount := getInputPrice(soldBaseAmt,
		result.BaseCoin.Amount, result.TokenCoin.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateBaseForSoldToken(soldToken sdk.Coin) (sdk.Int, error) {
	result, err := swap.QueryPool(soldToken.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	fee := sdk.MustNewDecFromStr(result.Fee)
	amount := getInputPrice(soldToken.Amount,
		result.TokenCoin.Amount, result.BaseCoin.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateTokenForSoldToken(boughtTokenDenom string,
	soldToken sdk.Coin) (sdk.Int, error) {
	if boughtTokenDenom == soldToken.Denom {
		return sdk.ZeroInt(), errors.New("invalid trade")
	}

	boughtBaseAmt, err := swap.EstimateBaseForSoldToken(soldToken)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return swap.EstimateTokenForSoldBase(boughtTokenDenom, boughtBaseAmt)
}

func (swap coinswapClient) EstimateTokenForBoughtBase(soldTokenDenom string,
	exactBoughtBaseAmt sdk.Int) (sdk.Int, error) {
	result, err := swap.QueryPool(soldTokenDenom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	fee := sdk.MustNewDecFromStr(result.Fee)
	amount := getOutputPrice(exactBoughtBaseAmt,
		result.TokenCoin.Amount, result.BaseCoin.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateBaseForBoughtToken(boughtToken sdk.Coin) (sdk.Int, error) {
	result, err := swap.QueryPool(boughtToken.Denom)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	fee := sdk.MustNewDecFromStr(result.Fee)
	amount := getOutputPrice(boughtToken.Amount,
		result.BaseCoin.Amount, result.TokenCoin.Amount, fee)
	return amount, nil
}

func (swap coinswapClient) EstimateTokenForBoughtToken(soldTokenDenom string,
	boughtToken sdk.Coin) (sdk.Int, error) {
	if soldTokenDenom == boughtToken.Denom {
		return sdk.ZeroInt(), errors.New("invalid trade")
	}

	soldBaseAmt, err := swap.EstimateBaseForBoughtToken(boughtToken)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return swap.EstimateTokenForBoughtBase(soldTokenDenom, soldBaseAmt)
}

func GetLiquidityDenomFrom(denom string) (string, error) {
	if denom == sdk.BaseDenom {
		return "", sdk.Wrapf("should not be base denom : %s", denom)
	}
	return fmt.Sprintf("swap%s", denom), nil
}

func GetTokenDenomFrom(liquidityDenom string) (string, error) {
	if !strings.HasPrefix(liquidityDenom, "swap") {
		return "", sdk.Wrapf("wrong liquidity denom : %s", liquidityDenom)
	}
	return strings.TrimPrefix(liquidityDenom, "swap"), nil
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
