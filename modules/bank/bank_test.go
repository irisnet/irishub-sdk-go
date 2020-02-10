package bank_test

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BankTestSuite struct {
	suite.Suite
	bank.Bank
	sender, validator string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(BankTestSuite))
}

func (tbs *BankTestSuite) SetupTest() {
	tc := sim.NewTestClient()
	tbs.Bank = tc.Bank
	tbs.sender = tc.GetTestSender()
	tbs.validator = tc.GetTestValidator()
}

func (tbs BankTestSuite) TestGetAccount() {
	acc, err := tbs.GetAccount("iaa1x3f572u057lv88mva2q3z40ls8pup9hsa74f9x")
	require.NoError(tbs.T(), err)
	fmt.Printf("%v", acc)
}

func (tbs BankTestSuite) TestGetTokenStats() {
	acc, err := tbs.GetTokenStats("iris")
	require.NoError(tbs.T(), err)
	fmt.Printf("%v", acc)
}

func (tbs BankTestSuite) TestSend() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	coins := types.NewCoins(coin)
	to := "iaa120v5ev44cwft687l0jcr5ec3vh2626vsschv7e"
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := tbs.Send(to, coins, baseTx)
	require.NoError(tbs.T(), err)
	require.True(tbs.T(), result.IsSuccess())
}

func (tbs BankTestSuite) TestBurn() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	coins := types.NewCoins(coin)
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
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
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := tbs.SetMemoRegexp("testMemo", baseTx)
	require.NoError(tbs.T(), err)
	require.True(tbs.T(), result.IsSuccess())

	acc, err := tbs.GetAccount(tbs.sender)
	require.NoError(tbs.T(), err)
	require.Equal(tbs.T(), "testMemo", acc.MemoRegexp)
}
