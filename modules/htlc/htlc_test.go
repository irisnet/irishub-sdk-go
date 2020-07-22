package htlc

import (
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

//type BankTestSuite struct {
//	suite.Suite
//	*test.MockClient
//}
//
//func TestKeeperTestSuite(t *testing.T) {
//	suite.Run(t, new(BankTestSuite))
//}
//
//func (bts *BankTestSuite) SetupTest() {
//	tc := test.GetMock()
//	bts.MockClient = tc
//}
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
	hts.Htlc().QueryHTLC()
}
