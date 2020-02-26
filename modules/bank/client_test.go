package bank_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/irisnet/irishub-sdk-go/types"
)

type BankTestSuite struct {
	suite.Suite
	sim.TestClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(BankTestSuite))
}

func (bts *BankTestSuite) SetupTest() {
	tc := sim.NewClient()
	bts.TestClient = tc
}

func (bts BankTestSuite) TestGetAccount() {
	acc, err := bts.QueryAccount("faa1d3mf696gvtwq2dfx03ghe64akf6t5vyz6pe3le")
	require.NoError(bts.T(), err)
	fmt.Printf("%v", acc)
}

func (bts BankTestSuite) TestGetTokenStats() {
	acc, err := bts.QueryTokenStats("iris")
	require.NoError(bts.T(), err)
	fmt.Printf("%v", acc)
}

func (bts BankTestSuite) TestSend() {
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

	toAccBefore, err := bts.QueryAccount(to)
	require.NoError(bts.T(), err)

	result, err := bts.Send(to, coins, baseTx)
	require.NoError(bts.T(), err)
	require.True(bts.T(), result.IsSuccess())

	toAccAfter, err := bts.QueryAccount(to)
	require.NoError(bts.T(), err)
	require.Equal(bts.T(),
		toAccBefore.Coins.Add(coins).String(),
		toAccAfter.GetCoins().String(),
	)
}

func (bts BankTestSuite) TestBurn() {
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
	result, err := bts.Burn(coins, baseTx)
	require.NoError(bts.T(), err)
	require.True(bts.T(), result.IsSuccess())
}

func (bts BankTestSuite) TestSetMemoRegexp() {
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := bts.SetMemoRegexp("testMemo", baseTx)
	require.NoError(bts.T(), err)
	require.True(bts.T(), result.IsSuccess())

	acc, err := bts.QueryAccount(bts.Sender().String())
	require.NoError(bts.T(), err)
	require.Equal(bts.T(), "testMemo", acc.MemoRegexp)
}
