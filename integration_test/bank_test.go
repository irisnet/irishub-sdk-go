package integration_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
)

func (s IntegrationTestSuite) TestAccount() {
	account, err := s.Bank.QueryAccount(s.Account().Address.String())
	require.NoError(s.T(), err)
	fmt.Print(account)
}
