package asset_test

import (
	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AssetTestSuite struct {
	suite.Suite
	test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(AssetTestSuite))
}

func (ats *AssetTestSuite) SetupTest() {
	tc := test.NewMockClient()
	ats.MockClient = tc
}

func (ats AssetTestSuite) TestQueryToken() {
	token, err := ats.Asset().QueryToken("iris")
	ats.NoError(err)
	ats.Equal(sdk.IRIS, token)
}

func (ats AssetTestSuite) TestQueryFees() {
	feeToken, err := ats.Asset().QueryFees("eth")
	ats.NoError(err)
	ats.Equal(false, feeToken.Exist)
	ats.Equal("60000000000000000000000iris-atto", feeToken.IssueFee.String())
	ats.Equal("6000000000000000000000iris-atto", feeToken.MintFee.String())
}
