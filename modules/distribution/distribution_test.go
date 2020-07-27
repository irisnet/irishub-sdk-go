package distribution_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/suite"
)

type DistrTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestDistrTestSuite(t *testing.T) {
	suite.Run(t, new(DistrTestSuite))
}

func (dts *DistrTestSuite) SetupTest() {
	tc := test.GetMock()
	dts.MockClient = tc
}

func (dts *DistrTestSuite) TestQueryRewards() {
	r, err := dts.Distr().QueryRewards("iaa18e2e9fxxrr88k78gg7fhuuqgccfv8self9ye65")
	require.NoError(dts.T(), err)
	require.NotEmpty(dts.T(), r)
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
	require.NoError(dts.T(), err)
	require.NotEmpty(dts.T(), rs.Hash)
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
	require.NoError(dts.T(), err)
	require.NotEmpty(dts.T(), rs.Hash)
}

func (dts *DistrTestSuite) TestQueryWithdrawAddr() {
	res, err := dts.Distr().QueryWithdrawAddr("iaa18e2e9fxxrr88k78gg7fhuuqgccfv8self9ye65")
	require.NoError(dts.T(), err)
	require.NotEmpty(dts.T(), res)
}
