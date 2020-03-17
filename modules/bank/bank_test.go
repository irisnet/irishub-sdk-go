package bank_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/irisnet/irishub-sdk-go/types"
)

type BankTestSuite struct {
	suite.Suite
	test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(BankTestSuite))
}

func (bts *BankTestSuite) SetupTest() {
	tc := test.NewMockClient()
	bts.MockClient = tc
}

func (bts BankTestSuite) TestGetAccount() {
	acc, err := bts.Bank().QueryAccount(bts.Account().Address.String())
	require.NoError(bts.T(), err)
	fmt.Printf("%v", acc)
}

func (bts BankTestSuite) TestGetTokenStats() {
	acc, err := bts.Bank().QueryTokenStats("iris")
	require.NoError(bts.T(), err)
	fmt.Printf("%v", acc)
}

func (bts BankTestSuite) TestSend() {
	coins, err := types.ParseDecCoins("0.1iris")
	require.NoError(bts.T(), err)
	to := "faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
	baseTx := types.BaseTx{
		From:     bts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: bts.Account().Password,
	}

	result, err := bts.Bank().Send(to, coins, baseTx)
	require.NoError(bts.T(), err)
	require.NotEmpty(bts.T(), result.Hash)
}

func (bts BankTestSuite) TestBurn() {
	amt, err := types.NewDecFromStr("0.1")
	require.NoError(bts.T(), err)
	coin := types.NewDecCoinFromDec("iris", amt)
	coins := types.NewDecCoins(coin)
	baseTx := types.BaseTx{
		From:     bts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: bts.Account().Password,
	}
	result, err := bts.Bank().Burn(coins, baseTx)
	require.NoError(bts.T(), err)
	require.NotEmpty(bts.T(), result.Hash)
}

func (bts BankTestSuite) TestSetMemoRegexp() {
	baseTx := types.BaseTx{
		From:     bts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: bts.Account().Password,
	}
	result, err := bts.Bank().SetMemoRegexp("testMemo", baseTx)
	require.NoError(bts.T(), err)
	require.NotEmpty(bts.T(), result.Hash)

	acc, err := bts.Bank().QueryAccount(bts.Account().Address.String())
	require.NoError(bts.T(), err)
	require.Equal(bts.T(), "testMemo", acc.MemoRegexp)
}
