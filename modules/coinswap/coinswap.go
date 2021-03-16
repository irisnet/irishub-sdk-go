package coinswap

import (
	"context"
	"fmt"
	"strings"

	"github.com/irisnet/irishub-sdk-go/codec"
	"github.com/irisnet/irishub-sdk-go/codec/types"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type coinswapClient struct {
	sdk.BaseClient
	codec.Marshaler
	queryTotalSupply func() (sdk.Coins, sdk.Error)
}

func NewClient(bc sdk.BaseClient, cdc codec.Marshaler,queryTotalSupply func() (sdk.Coins, sdk.Error)) Client {
	return coinswapClient{
		BaseClient: bc,
		Marshaler:  cdc,
		queryTotalSupply: queryTotalSupply,
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

	liquidityDenom := getLiquidityDenomFrom(request.MaxToken.Denom)
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

	tokenDenom, er := getTokenDenomFrom(request.Liquidity.Denom)
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

func (swap coinswapClient) QueryPool(denom string) (*QueryPoolResponse, error) {
	conn, err := swap.GenConn()
	defer func() { _ = conn.Close() }()
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

func (swap coinswapClient) QueryAllPools() (*QueryAllPoolsResponse,error){
	coins,err :=swap.queryTotalSupply()
	if err != nil {
		return nil, sdk.Wrap(err)
	}

	var pools []QueryPoolResponse
	for _,coin := range coins {
		denom,err := getTokenDenomFrom(coin.Denom)
		if err != nil {
			continue
		}
		res,err := swap.QueryPool(denom)
		if err != nil {
			return nil, sdk.Wrap(err)
		}
		pools= append(pools,*res)
	}
	return &QueryAllPoolsResponse{pools}, err
}

func getLiquidityDenomFrom(denom string) string {
	return fmt.Sprintf("swap/%s", denom)
}

func getTokenDenomFrom(liquidityDenom string) (string, error) {
	sp := strings.Split(liquidityDenom, "/")
	if len(sp) != 2 {
		return "", sdk.Wrapf("wrong liquidity denom : %s", liquidityDenom)
	}
	return sp[1], nil
}
