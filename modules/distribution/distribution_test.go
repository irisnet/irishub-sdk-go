package distribution_test

import (
	"testing"

	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/suite"
)

type DistrTestSuite struct {
	suite.Suite
	test.MockClient
}

func TestDistrTestSuite(t *testing.T) {
	suite.Run(t, new(DistrTestSuite))
}

func (dts *DistrTestSuite) SetupTest() {
	tc := test.NewMockClient()
	dts.MockClient = tc
}

func (dts *DistrTestSuite) TestQueryRewards() {
	r, err := dts.Distr().QueryRewards(dts.Account().Address.String())
	dts.NoError(err)
	dts.NotEmpty(r)
}

func (dts *DistrTestSuite) TestSetWithdrawAddr() {
	baseTx := sdk.BaseTx{
		From:     dts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: dts.Account().Password,
	}

	rs, err := dts.Distr().SetWithdrawAddr(dts.Account().Address.String(), baseTx)
	dts.NoError(err)
	dts.NotEmpty(rs.Hash)
}

func (dts *DistrTestSuite) TestWithdrawRewards() {
	baseTx := sdk.BaseTx{
		From:     dts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: dts.Account().Password,
	}

	rs, err := dts.Distr().WithdrawRewards(true, "", baseTx)
	dts.NoError(err)
	dts.NotEmpty(rs.Hash)
}
