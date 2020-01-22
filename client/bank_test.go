package client_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func (c *ClientTestSuite) TestGetAccount() {
	acc, err := c.GetAccount("iaa1x3f572u057lv88mva2q3z40ls8pup9hsa74f9x")
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestGetTokenStats() {
	acc, err := c.GetTokenStats("iris")
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestSend() {
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
	result, err := c.Send(to, coins, baseTx)
	require.NoError(c.T(), err)
	require.True(c.T(), result.IsSuccess())
}

func (c *ClientTestSuite) TestBurn() {
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
	result, err := c.Burn(coins, baseTx)
	require.NoError(c.T(), err)
	require.True(c.T(), result.IsSuccess())
}

func (c *ClientTestSuite) TestSetMemoRegexp() {
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	result, err := c.SetMemoRegexp("testMemo", baseTx)
	require.NoError(c.T(), err)
	require.True(c.T(), result.IsSuccess())

	acc, err := c.GetAccount(addr)
	require.NoError(c.T(), err)
	require.Equal(c.T(), "testMemo", acc.MemoRegexp)
}
