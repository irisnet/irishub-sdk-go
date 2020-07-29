package params_test

import (
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ParamTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestDistrTestSuite(t *testing.T) {
	suite.Run(t, new(ParamTestSuite))
}

func (pts *ParamTestSuite) SetupTest() {
	tc := test.GetMock()
	pts.MockClient = tc
}

func (pts *ParamTestSuite) TestQueryParam() {
	subspace := "staking"
	key := "MaxValidators"
	params, err := pts.Params().QueryParams(subspace, key)
	require.NoError(pts.T(), err)
	require.Equal(pts.T(), params.Subspace, subspace)
	require.Equal(pts.T(), params.Key, key)
}
