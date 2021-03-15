package integration_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/coinswap"
	"time"

	"github.com/stretchr/testify/require"

	sdk "github.com/irisnet/irishub-sdk-go/types"
)

func (s IntegrationTestSuite) TestCoinSwap() {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	request := coinswap.AddLiquidityRequest{
		MaxToken:     sdk.Coin{
			Denom: "ubnb",
			Amount: sdk.NewInt(1000_000_000),
		},
		BaseAmt:      sdk.NewInt(1000_000_000),
		MinLiquidity: sdk.NewInt(1),
		Deadline:     time.Now().Add(time.Hour).Unix(),
	}

	res,err := s.Swap.AddLiquidity(request,baseTx)
	require.NoError(s.T(), err)
	fmt.Println(res.TxHash)
	fmt.Println(res.Liquidity.String())
	fmt.Println(res.TokenAmt.String())
}

func (s IntegrationTestSuite) TestSwapCoin()  {
	baseTx := sdk.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: s.Account().Password,
	}

	swapReq := coinswap.SwapCoinRequest{
		Input:     sdk.Coin{
			Denom: "uiris",
			Amount: sdk.NewInt(2_000_000),
		},
		Output:      sdk.Coin{
			Denom: "ubnb",
			Amount: sdk.NewInt(1_000_000),
		},
		Deadline:     time.Now().Add(time.Hour).Unix(),
		IsBuyOrder:   true,
	}

	res1,err := s.Swap.SwapCoin(swapReq,baseTx)
	require.NoError(s.T(), err)
	fmt.Println(res1.TxHash)
	fmt.Println(res1.InputAmt.String())
	fmt.Println(res1.OutputAmt.String())
}

func (s IntegrationTestSuite) TestQueryPool()  {
	res,err := s.Swap.QueryPool("ubnb")
	require.NoError(s.T(), err)
	fmt.Println(res.Liquidity.String())
	fmt.Println(res.TokenCoin.String())
	fmt.Println(res.BaseCoin.String())
}
