package integration_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub-sdk-go/modules/coinswap"
	"github.com/irisnet/irishub-sdk-go/modules/token"
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

	request := coinswap.AddLiquidityRequest{
		MaxToken: sdk.Coin{
			Denom:  "ubnb",
			Amount: sdk.NewInt(1000_000_000),
		},
		BaseAmt:      sdk.NewInt(1000_000_000),
		MinLiquidity: sdk.NewInt(1000_000_000),
		Deadline:     time.Now().Add(time.Hour).Unix(),
	}

	res, err := s.Swap.AddLiquidity(request, baseTx)
	require.NoError(s.T(), err)
	require.True(s.T(), res.Liquidity.GTE(request.MinLiquidity))
	require.NotEmpty(s.T(), res.TxHash)
	require.True(s.T(), request.MaxToken.Amount.GTE(res.TokenAmt))

	boughtCoin := sdk.NewCoin("uiris", sdk.NewInt(100))
	deadline := time.Now().Add(10 * time.Second).Unix()
	resp, err := s.Swap.BuyTokenWithAutoEstimate("ubnb", boughtCoin, deadline, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.TxHash)
	require.True(s.T(), resp.InputAmt.Equal(sdk.NewInt(101)))

	soldCoin := sdk.NewCoin("uiris", sdk.NewInt(100))
	resp, err = s.Swap.SellTokenWithAutoEstimate("ubnb", soldCoin, deadline, baseTx)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resp.TxHash)
	require.True(s.T(), resp.OutputAmt.Equal(sdk.NewInt(99)))
}

func (s IntegrationTestSuite) TestQuery() {
	res, err := s.Swap.QueryAllPools(sdk.PageRequest{
		Offset:     0,
		Limit:      10,
		CountTotal: false,
	})
	require.NoError(s.T(), err)
	fmt.Println(res)

	res1, err1 := s.Swap.QueryPool("lpt-1")
	require.NoError(s.T(), err1)
	fmt.Println(res1)

}
