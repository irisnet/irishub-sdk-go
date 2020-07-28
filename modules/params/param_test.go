package params_test

import (
	"fmt"
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
	params, err := pts.Params().QueryParams()
	fmt.Println(params)
	require.NoError(pts.T(), err)
}
