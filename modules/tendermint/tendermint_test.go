package tendermint_test

import (
	"fmt"
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TendermintTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TendermintTestSuite))
}

func (tts *TendermintTestSuite) SetupTest() {
	tc := test.GetMock()
	tts.MockClient = tc
}

func (tts *TendermintTestSuite) TestQueryQueryBlock() {
	var height int64 = 34433
	block, err := tts.Tendermint().QueryBlock(height)
	require.NoError(tts.T(), err)
	require.Equal(tts.T(), height, block.Height)
}

func (tts *TendermintTestSuite) TestQueryBlockLatest() {
	block, err := tts.Tendermint().QueryBlockLatest()
	require.NoError(tts.T(), err)
	require.NotEmpty(tts.T(), block)
}

func (tts *TendermintTestSuite) TestQueryBlockResult() {
	result, err := tts.Tendermint().QueryBlockResult(34433)
	fmt.Println(result)
	require.NoError(tts.T(), err)
	require.Equal(tts.T(), int64(1443), result.Height)
}

func (tts *TendermintTestSuite) TestQueryTx() {
	tx, err := tts.Tendermint().QueryTx("B664260AD9E9E8400B4B865123C84333A4974194ED92BEA185EDFFD9BECEF5D7")
	fmt.Println(tx)
	require.NoError(tts.T(), err)
	fmt.Println(tx)

}

func (tts *TendermintTestSuite) TestQueryValidators() {
	result, err := tts.Tendermint().QueryValidators(1007)
	require.NoError(tts.T(), err)
	require.Len(tts.T(), result.Validators, 1)
}

func (tts *TendermintTestSuite) TestQueryValidatorsLatest() {
	result, err := tts.Tendermint().QueryValidatorsLatest()
	require.NoError(tts.T(), err)
	require.Len(tts.T(), result.Validators, 1)
}

func (tts *TendermintTestSuite) TestQueryNodeInfo() {
	result, err := tts.Tendermint().QueryNodeInfo()
	require.NoError(tts.T(), err)
	require.NotEmpty(tts.T(), result)
}

func (tts *TendermintTestSuite) TestQueryNodeVersion() {
	result, err := tts.Tendermint().QueryNodeVersion()
	require.NoError(tts.T(), err)
	require.NotEmpty(tts.T(), result)
}

func (tts *TendermintTestSuite) TestQueryGenesis() {
	result, err := tts.Tendermint().QueryGenesis()
	require.NoError(tts.T(), err)
	require.NotEmpty(tts.T(), result)
}
