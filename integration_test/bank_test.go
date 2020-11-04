package integration_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/types"
)

func (s IntegrationTestSuite) TestAccount() {
	account, err := s.Bank.QueryAccount(s.Account().Address.String())
	s.NoError(err)
	fmt.Print(account)
}

func (s IntegrationTestSuite) TestSend() {
	coins, err := types.ParseDecCoins("1iris")
	s.NoError(err)
	to := s.GetRandAccount().Address.String()
	baseTx := types.BaseTx{
		From:     s.Account().Name,
		Gas:      200000,
		Fee:      coins,
		Memo:     "TEST",
		Mode:     types.Commit,
		Password: s.Account().Password,
	}

	res, err := s.Bank.Send(to, coins, baseTx)
	s.NoError(err)
	fmt.Println(res)
}
