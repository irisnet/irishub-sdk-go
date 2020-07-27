package asset_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/test"
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

func (ats *AssetTestSuite) TestQueryTokens() {
	token, err := ats.Asset().QueryTokens()
	require.NoError(ats.T(), err)
	require.NotEmpty(ats.T(), token)
}

func (ats AssetTestSuite) TestQueryToken() {
	token, err := ats.Asset().QueryTokenDenom("stake")
	if err != nil {
		ats.Error(err)
	}
	fmt.Println(token)
	//fmt.Println(err.Error())
	require.NoError(ats.T(), err)
}
