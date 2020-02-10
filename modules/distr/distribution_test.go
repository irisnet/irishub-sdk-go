package distr_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/modules/distr"
	"github.com/irisnet/irishub-sdk-go/types"
	"testing"

	"github.com/irisnet/irishub-sdk-go/sim"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DistrTestSuite struct {
	types.Codec
	suite.Suite
	distr.Distr
	sender, validator string
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(DistrTestSuite))
}

func (gts *DistrTestSuite) SetupTest() {
	tc := sim.NewTestClient()
	gts.Distr = tc.Distr
	gts.sender = tc.GetTestSender()
}

func (gts DistrTestSuite) TestQueryRewards() {
	deposit, err := gts.QueryRewards(gts.sender)
	require.NoError(gts.T(), err)
	fmt.Printf("%v", deposit)
}
