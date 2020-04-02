package bank_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"sync"
	"testing"
	"time"
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
	bts.NoError(err)
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
	amt, err := types.NewDecimalFromStr("0.1")
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

	var du time.Duration
	fmt.Println(du.Nanoseconds())
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

func (bts BankTestSuite) TestMultiSend() {
	baseTx := types.BaseTx{
		From:     bts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     types.Commit,
		Password: bts.Account().Password,
	}

	coins, e := types.ParseDecCoins("1000iris")
	require.NoError(bts.T(), e)

	bank := bts.Bank()

	var accNum = 11
	var acc = make([]string, accNum)
	var receipts = make([]rpc.Receipt, accNum)
	for i := 0; i < accNum; i++ {
		acc[i] = fmt.Sprintf("%s%d", "testBank", i)
		addr, _, err := bts.Keys().Add(acc[i], "1234567890")

		require.NoError(bts.T(), err)
		require.NotEmpty(bts.T(), addr)

		receipts[i] = rpc.Receipt{
			Address: addr,
			Amount:  coins,
		}
	}

	_, err := bank.MultiSend(receipts, baseTx)
	require.NoError(bts.T(), err)

	coins, e = types.ParseDecCoins("1iris")
	require.NoError(bts.T(), e)

	to := "faa1leqs0fg0nsav2u3vdt4gat0mpascl52kl2huel"

	begin := time.Now()
	var wait sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wait.Add(1)
		index := rand.Intn(accNum)
		go func() {
			defer wait.Done()
			_, err := bank.Send(to, coins, types.BaseTx{
				From:     acc[index],
				Gas:      20000,
				Memo:     "test",
				Mode:     types.Async,
				Password: "1234567890",
			})
			require.NoError(bts.T(), err)
		}()
	}
	wait.Wait()
	end := time.Now()
	fmt.Println(fmt.Sprintf("total senconds:%s", end.Sub(begin).String()))
}
