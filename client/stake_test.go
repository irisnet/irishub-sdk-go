package client_test

import (
	"encoding/json"
	"fmt"

	"github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
)

func (c *ClientTestSuite) TestQueryDelegation() {
	acc, err := c.QueryDelegation(addr, valAddr)
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestQueryDelegations() {
	acc, err := c.QueryDelegations(addr)
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestQueryValidator() {
	acc, err := c.QueryValidator(valAddr)
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestQueryAllValidators() {
	val, err := c.QueryAllValidators()
	require.NoError(c.T(), err)
	bz, err := json.Marshal(val)
	require.NoError(c.T(), err)
	fmt.Printf(string(bz))
}

func (c *ClientTestSuite) TestQueryValidators() {
	acc, err := c.QueryValidators(1, 100)
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestDelegate() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	acc, err := c.Delegate(valAddr, coin, baseTx)
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestUnDelegate() {
	amt := types.NewIntWithDecimal(1, 18)
	coin := types.NewCoin("iris-atto", amt)
	baseTx := types.BaseTx{
		From: "test1",
		Gas:  "20000",
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: types.Commit,
	}
	acc, err := c.Undelegate(valAddr, coin, baseTx)
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}
