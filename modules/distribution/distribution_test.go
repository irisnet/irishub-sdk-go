package distribution_test

import (
	"testing"

	sdk "github.com/irisnet/irishub-sdk-go/types"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
)

type DistrTestSuite struct {
	suite.Suite
	test.TestClient
}

func TestDistrTestSuite(t *testing.T) {
	suite.Run(t, new(DistrTestSuite))
}

func (dts *DistrTestSuite) SetupTest() {
	tc := test.NewClient()
	dts.TestClient = tc
}

func (dts *DistrTestSuite) TestQueryRewards() {
	r, err := dts.Distr().QueryRewards(dts.Sender().String())
	require.NoError(dts.T(), err)
	require.NotEmpty(dts.T(), r)
}

func (dts *DistrTestSuite) TestSetWithdrawAddr() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Memo: "test",
		Mode: sdk.Commit,
	}

	rs, err := dts.Distr().SetWithdrawAddr(dts.Sender().String(), baseTx)
	require.NoError(dts.T(), err)
	require.NotEmpty(dts.T(), rs.Hash)
}

func (dts *DistrTestSuite) TestWithdrawRewards() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Memo: "test",
		Mode: sdk.Commit,
	}

	rs, err := dts.Distr().WithdrawRewards(true, "", baseTx)
	require.NoError(dts.T(), err)
	require.NotEmpty(dts.T(), rs.Hash)
}
