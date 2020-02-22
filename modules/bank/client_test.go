package bank_test

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BankTestSuite struct {
	suite.Suite
	types.Bank
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(BankTestSuite))
}

func (tbs *BankTestSuite) SetupTest() {
	tc := sim.NewClient()
	tbs.Bank = tc.Bank
}

func (tbs BankTestSuite) TestGetAccount() {
	acc, err := tbs.QueryAccount("faa1d3mf696gvtwq2dfx03ghe64akf6t5vyz6pe3le")
	require.NoError(tbs.T(), err)
	fmt.Printf("%v", acc)
}

func (tbs BankTestSuite) TestGetTokenStats() {
	acc, err := tbs.QueryTokenStats("iris")
	require.NoError(tbs.T(), err)
	fmt.Printf("%v", acc)
}

func (tbs BankTestSuite) TestSend() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	coins := types.NewCoins(coin)
	to := "faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}

	toAccBefore, err := tbs.QueryAccount(to)
	require.NoError(tbs.T(), err)

	result, err := tbs.Send(to, coins, baseTx)
	require.NoError(tbs.T(), err)
	require.True(tbs.T(), result.IsSuccess())

	toAccAfter, err := tbs.QueryAccount(to)
	require.NoError(tbs.T(), err)
	require.Equal(tbs.T(),
		toAccBefore.Coins.Add(coins).String(),
		toAccAfter.GetCoins().String(),
	)
}

func (tbs BankTestSuite) TestBurn() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	coins := types.NewCoins(coin)
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := tbs.Burn(coins, baseTx)
	require.NoError(tbs.T(), err)
	require.True(tbs.T(), result.IsSuccess())
}

func (tbs BankTestSuite) TestSetMemoRegexp() {
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := tbs.SetMemoRegexp("testMemo", baseTx)
	require.NoError(tbs.T(), err)
	require.True(tbs.T(), result.IsSuccess())

	acc, err := tbs.QueryAccount(sim.Addr)
	require.NoError(tbs.T(), err)
	require.Equal(tbs.T(), "testMemo", acc.MemoRegexp)
}
