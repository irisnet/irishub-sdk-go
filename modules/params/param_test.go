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

func (pts *ParamTestSuite) TestQueryParamsBySubAndKey() {
	subspace := "staking"
	key := "MaxValidators"
	params, err := pts.Params().QueryParamsBySubAndKey(subspace, key)
	require.NoError(pts.T(), err)
	require.Equal(pts.T(), params.Subspace, subspace)
	require.Equal(pts.T(), params.Key, key)
}

func (pts *ParamTestSuite) TestQueryParams() {
	params, err := pts.Params().QueryParams("token")
	require.NoError(pts.T(), err)
	require.NotEmpty(pts.T(), params)
}
