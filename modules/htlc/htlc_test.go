package htlc_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HtlcTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(HtlcTestSuite))
}

func (hts *HtlcTestSuite) SetupTest() {
	tc := test.GetMock()
	hts.MockClient = tc
}

func (hts HtlcTestSuite) TestGetTokenStats() {
	htlc, err := hts.Htlc().QueryHtlc("123")
	fmt.Println(htlc)
	require.NoError(hts.T(), err)
}
