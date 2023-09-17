package integration_test

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/irisnet/irishub-sdk-go/modules/coinswap"
	"github.com/irisnet/irishub-sdk-go/modules/token"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func (s IntegrationTestSuite) TestCoinSwap() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Fee:	  sdk.NewDecCoins(sdk.NewDecCoin("uiris", sdk.NewInt(10))),
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	bnb, err := s.Token.QueryToken("bnb")
	if bnb.Symbol == "" {
		issueTokenReq := token.IssueTokenRequest{
			Symbol:        "bnb",
			Name:          s.RandStringOfLength(8),
			Scale:         6,
			MinUnit:       "ubnb",
			InitialSupply: 10000000,
			MaxSupply:     21000000,
			Mintable:      true,
		}
		result, er := s.Token.IssueToken(issueTokenReq, baseTx)
		require.NoError(s.T(), er)
		require.NotEmpty(s.T(), result.Hash)
	}
	busd, err := s.Token.QueryToken("busd")
	if busd.Symbol == "" {
		issueTokenReq := token.IssueTokenRequest{
			Symbol:        "busd",
			Name:          s.RandStringOfLength(8),
			Scale:         6,
			MinUnit:       "ubusd",
			InitialSupply: 10000000,
			MaxSupply:     21000000,
			Mintable:      true,
		}
		time.Sleep(5*time.Second)
		result, er := s.Token.IssueToken(issueTokenReq, baseTx)
		require.NoError(s.T(), er)
		require.NotEmpty(s.T(), result.Hash)
	}

	request := coinswap.AddLiquidityRequest{
		MaxToken: sdk.Coin{
			Denom:  "ubnb",
			Amount: sdk.NewInt(1000_000_0000),
		},
		BaseAmt:      sdk.NewInt(1000_000_000),
		MinLiquidity: sdk.NewInt(1000_000_000),
		Deadline:     time.Now().Add(time.Hour).Unix(),
	}
	time.Sleep(5*time.Second)
	res, err := s.Swap.AddLiquidity(request, baseTx)
	fmt.Println(">>>>> add bnb", res, err)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), res.TxHash)
	require.True(s.T(), res.Liquidity.Amount.GTE(request.MinLiquidity))
	require.True(s.T(), request.MaxToken.Amount.GTE(res.TokenAmt))
	
	time.Sleep(5*time.Second)
	p1, err := s.Swap.QueryPool(res.Liquidity.Denom)
	require.NoError(s.T(), err)
	require.True(s.T(), p1.Pool.Token.Denom == request.MaxToken.Denom)

	rlr := coinswap.RemoveLiquidityRequest{
		MinTokenAmt: sdk.NewInt(100000000),
		MinBaseAmt: sdk.NewInt(100000000),
		Liquidity: sdk.NewCoin(res.Liquidity.Denom, sdk.NewInt(500000000)),
		Deadline: time.Now().Add(10 * time.Minute).Unix(),
	}
	res2, err := s.Swap.RemoveLiquidity(rlr, baseTx)
	fmt.Println(">>>>> remove bnb", p1, rlr, res2)
	require.NoError(s.T(), err)
	require.True(s.T(), res2.TokenAmt.Equal(p1.Pool.Token.Amount.Mul(rlr.Liquidity.Amount).Quo(p1.Pool.Lpt.Amount)))
	
	boughtCoin := sdk.NewCoin("uiris", sdk.NewInt(100))
	deadline := time.Now().Add(1 * time.Hour).Unix()
	time.Sleep(5*time.Second)
	p1, _ = s.Swap.QueryPool(res.Liquidity.Denom)
	resp, err := s.Swap.BuyTokenWithAutoEstimate("ubnb", boughtCoin, []sdk.PoolInfo{p1.Pool}, deadline, baseTx)
	fee, _ := strconv.ParseFloat(p1.Pool.Fee, 64)
	result := CS(-float64(boughtCoin.Amount.Uint64()), float64(p1.Pool.Standard.Amount.Int64()), float64(p1.Pool.Token.Amount.Int64()), fee)
	fmt.Println(">>>>> swap bnb", p1, resp, resp.InputAmt, result)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.TxHash)
	require.True(s.T(), math.Abs(float64(resp.InputAmt.Int64()) - result) < 1)

	estimatedAmt, err := s.Swap.EstimateBaseForBoughtToken(sdk.NewCoin("ubnb", sdk.NewInt(100000)), &p1.Pool)
	result = CS(-float64(100000), float64(p1.Pool.Token.Amount.Int64()), float64(p1.Pool.Standard.Amount.Int64()), fee)
	require.NoError(s.T(), err)
	fmt.Println(">>>>> EstimateBaseForBoughtToken", estimatedAmt, result)
	require.True(s.T(), math.Abs(float64(estimatedAmt.Int64()) - result) < 1)

	estimatedAmt, err = s.Swap.EstimateBaseForSoldToken(sdk.NewCoin("ubnb", sdk.NewInt(100000)), &p1.Pool)
	require.NoError(s.T(), err)
	result = CS(float64(100000), float64(p1.Pool.Token.Amount.Int64()), float64(p1.Pool.Standard.Amount.Int64()), fee)
	fmt.Println(">>>>> EstimateBaseForSoldToken", estimatedAmt, result)
	require.True(s.T(), math.Abs(float64(estimatedAmt.Int64()) - result) < 1)
	
	estimatedAmt, err = s.Swap.EstimateTokenForBoughtBase("ubnb", sdk.NewInt(1000000), &p1.Pool)
	result = CS(-float64(1000000), float64(p1.Pool.Standard.Amount.Int64()), float64(p1.Pool.Token.Amount.Int64()), fee)
	fmt.Println(">>>>> EstimateTokenForBoughtBase", estimatedAmt, result)
	require.NoError(s.T(), err)
	require.True(s.T(), math.Abs(float64(estimatedAmt.Int64()) - result) < 1)
	
	estimatedAmt, err = s.Swap.EstimateTokenForSoldBase("ubnb", sdk.NewInt(1000000), &p1.Pool)
	result = CS(float64(1000000), float64(p1.Pool.Standard.Amount.Int64()), float64(p1.Pool.Token.Amount.Int64()), fee)
	fmt.Println(">>>>> EstimateTokenForSoldBase", estimatedAmt, result)
	require.NoError(s.T(), err)
	require.True(s.T(), math.Abs(float64(estimatedAmt.Int64()) - result) < 1)
	

	request = coinswap.AddLiquidityRequest{
		MaxToken: sdk.Coin{
			Denom:  "ubusd",
			Amount: sdk.NewInt(1000_000_00),
		},
		BaseAmt:      sdk.NewInt(1000_000),
		MinLiquidity: sdk.NewInt(800_000),
		Deadline:     time.Now().Add(time.Hour).Unix(),
	}
	time.Sleep(5*time.Second)
	res, err = s.Swap.AddLiquidity(request, baseTx)
	fmt.Println(">>>>> add busd", res, err)
	require.NoError(s.T(), err)
	require.True(s.T(), res.Liquidity.Amount.GTE(request.MinLiquidity))
	require.NotEmpty(s.T(), res.TxHash)
	require.True(s.T(), request.MaxToken.Amount.GTE(res.TokenAmt))

	time.Sleep(5*time.Second)
	p2, err := s.Swap.QueryPool(res.Liquidity.Denom)
	fmt.Println(">>>>> query pool", p2, err, res.Liquidity.Denom)
	require.NoError(s.T(), err)
	require.True(s.T(), p2.Pool.Token.Denom == request.MaxToken.Denom)
	
	rlr = coinswap.RemoveLiquidityRequest{
		MinTokenAmt: sdk.NewInt(10000),
		MinBaseAmt: sdk.NewInt(10000),
		Liquidity: sdk.NewCoin(res.Liquidity.Denom, sdk.NewInt(50000)),
		Deadline: time.Now().Add(10 * time.Minute).Unix(),
	}
	res2, err = s.Swap.RemoveLiquidity(rlr, baseTx)
	fmt.Println(">>>>> remove busd", p2, rlr, res2)
	require.NoError(s.T(), err)
	require.True(s.T(), res2.TokenAmt.Equal(p2.Pool.Token.Amount.Mul(rlr.Liquidity.Amount).Quo(p2.Pool.Lpt.Amount)))

	soldCoin := sdk.NewCoin("uiris", sdk.NewInt(100))
	time.Sleep(5*time.Second)
	p2, _ = s.Swap.QueryPool(res.Liquidity.Denom)
	resp, err = s.Swap.SellTokenWithAutoEstimate("ubusd", soldCoin, []sdk.PoolInfo{p2.Pool}, deadline, baseTx)
	result = CS(float64(soldCoin.Amount.Int64()), float64(p2.Pool.Standard.Amount.Int64()), float64(p2.Pool.Token.Amount.Int64()), fee)
	fmt.Println(">>>>> swap busd", p2, resp, resp.OutputAmt, result)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.TxHash)
	fee, _ = strconv.ParseFloat(p2.Pool.Fee, 64)
	require.True(s.T(), math.Abs(float64(resp.OutputAmt.Int64()) - result) < 1)

	estimatedAmt, err = s.Swap.EstimateTokenForSoldBase("ubnb", sdk.NewInt(1000000), &p2.Pool)
	result = CS(float64(1000000), float64(p1.Pool.Standard.Amount.Int64()), float64(p1.Pool.Token.Amount.Int64()), fee)
	require.Error(s.T(), err)
	fmt.Println(">>>>> EstimateTokenForSoldBase", err, estimatedAmt, result)

	estimatedAmt, err = s.Swap.EstimateTokenForBoughtToken("ubnb", sdk.NewCoin("ubusd", sdk.NewInt(1000000)), []sdk.PoolInfo{p1.Pool, p2.Pool})
	result = CS2(-1000000, float64(p2.Pool.Token.Amount.Int64()), float64(p2.Pool.Standard.Amount.Int64()), float64(p1.Pool.Standard.Amount.Int64()), float64(p1.Pool.Token.Amount.Int64()), fee, fee)
	fmt.Println(">>>>> EstimateTokenForBoughtToken", estimatedAmt, result)
	require.NoError(s.T(), err)
	require.True(s.T(), math.Abs(float64(estimatedAmt.Int64()) - result) < 1)

	estimatedAmt, err = s.Swap.EstimateTokenForSoldToken("ubnb", sdk.NewCoin("ubusd", sdk.NewInt(10000000)), []sdk.PoolInfo{p1.Pool, p2.Pool})
	result = CS2(10000000, float64(p2.Pool.Token.Amount.Int64()), float64(p2.Pool.Standard.Amount.Int64()), float64(p1.Pool.Standard.Amount.Int64()), float64(p1.Pool.Token.Amount.Int64()), fee, fee)
	fmt.Println(">>>>> EstimateTokenForSoldToken", estimatedAmt, result)
	require.NoError(s.T(), err)
	require.True(s.T(), math.Abs(float64(estimatedAmt.Int64()) - result) < 1)

	estimatedAmt, err = s.Swap.EstimateTokenForSoldToken("ubnb", sdk.NewCoin("ubusd", sdk.NewInt(10000000)), []sdk.PoolInfo{p1.Pool, p1.Pool})
	result = CS2(10000000, float64(p2.Pool.Token.Amount.Int64()), float64(p2.Pool.Standard.Amount.Int64()), float64(p1.Pool.Standard.Amount.Int64()), float64(p1.Pool.Token.Amount.Int64()), fee, fee)
	require.Error(s.T(), err)
	fmt.Println(">>>>> EstimateTokenForSoldToken", err, estimatedAmt, result)

}

func CS(x, X, Y, f float64) float64 {
	if x > 0 {
		return x * (1-f) * Y / (x * (1-f) + X)
	}
	return -x * Y / (x + X) / (1-f)
}
func CS2(x, X, Z1, Z2, Y, f1, f2 float64) float64 {
	if x > 0 {
		return x*Y*Z1*(1-f1)*(1-f2) / (x*(1-f1)*(1-f2)*Z1 + (x*(1-f1)+X)*Z2)
	}
	return -x*Y*Z1 / (x*(1-f2)*Z1 + (x+X)*(1-f1)*(1-f2)*Z2)
}

func (s IntegrationTestSuite) TestQuery() {
	res, err := s.Swap.QueryAllPools(sdk.PageRequest{
		Offset:     0,
		Limit:      10,
		CountTotal: false,
	})
	require.NoError(s.T(), err)
	fmt.Println(res)

	p, err1 := s.Swap.QueryPool("lpt-5")
	require.NoError(s.T(), err1)
	fmt.Println(p, CS(-100000, float64(p.Pool.Token.Amount.Int64()), float64(p.Pool.Standard.Amount.Int64()), 0.003))
	fmt.Println(s.Swap.EstimateBaseForBoughtToken(sdk.NewCoin("ubnb", sdk.NewInt(100000)), &p.Pool))
}
