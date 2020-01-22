package client_test

import (
	"fmt"
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
	acc, err := c.QueryAllValidators()
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}

func (c *ClientTestSuite) TestQueryValidators() {
	acc, err := c.QueryValidators(1, 100)
	require.NoError(c.T(), err)
	fmt.Printf("%v", acc)
}
