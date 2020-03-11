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
	acc, err := bts.Bank().QueryAccount(bts.Sender().String())
	require.True(bts.T(), err.IsNil())
	fmt.Printf("%v", acc)
}

func (bts BankTestSuite) TestGetTokenStats() {
	acc, err := bts.Bank().QueryTokenStats("iris")
	require.True(bts.T(), err.IsNil())
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
		Memo: "test",
		Mode: types.Commit,
	}

	beforeCoin := types.NewCoins()
	toAccBefore, err := bts.Bank().QueryAccount(to)
	if err.IsNil() {
		beforeCoin = toAccBefore.GetCoins()
	}

	result, err := bts.Bank().Send(to, coins, baseTx)
	require.True(bts.T(), err.IsNil())
	require.NotEmpty(bts.T(), result.Hash)

	toAccAfter, err := bts.Bank().QueryAccount(to)
	require.True(bts.T(), err.IsNil())
	require.Equal(bts.T(),
		beforeCoin.Add(coins).String(),
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
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := bts.Bank().Burn(coins, baseTx)
	require.True(bts.T(), err.IsNil())
	require.NotEmpty(bts.T(), result.Hash)
}

func (bts BankTestSuite) TestSetMemoRegexp() {
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  20000,
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := bts.Bank().SetMemoRegexp("testMemo", baseTx)
	require.True(bts.T(), err.IsNil())
	require.NotEmpty(bts.T(), result.Hash)

	acc, err := bts.Bank().QueryAccount(bts.Sender().String())
	require.True(bts.T(), err.IsNil())
	require.Equal(bts.T(), "testMemo", acc.MemoRegexp)
}
