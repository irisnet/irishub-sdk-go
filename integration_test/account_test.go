package integration_test

import "github.com/stretchr/testify/require"

func (s IntegrationTestSuite) TestAccount() {
	_, err := s.QueryAccount(s.Account().Address.String())
	require.NoError(s.T(), err)

}
