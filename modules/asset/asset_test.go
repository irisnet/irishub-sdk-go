package asset_test

import (
	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AssetTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(AssetTestSuite))
}

func (ats *AssetTestSuite) SetupTest() {
	tc := test.GetMock()
	ats.MockClient = tc
}

func (ats AssetTestSuite) TestQueryToken() {
	token, err := ats.Asset().QueryToken("iris")
	require.NoError(ats.T(), err)
	require.Equal(ats.T(), sdk.IRIS, token)
}

func (ats AssetTestSuite) TestQueryFees() {
	feeToken, err := ats.Asset().QueryFees("eth")
	require.NoError(ats.T(), err)
	require.Equal(ats.T(), false, feeToken.Exist)
	require.Equal(ats.T(), "60000000000000000000000iris-atto", feeToken.IssueFee.String())
	require.Equal(ats.T(), "6000000000000000000000iris-atto", feeToken.MintFee.String())
}
