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
	var height int64 = 1006
	block, err := tts.Tendermint().QueryBlock(height)
	require.NoError(tts.T(), err)
	require.Equal(tts.T(), height, block.Height)
}

func (tts *TendermintTestSuite) TestQueryQueryLatest() {
	block, err := tts.Tendermint().QueryBlockLatest()
	require.NoError(tts.T(), err)
	require.NotEmpty(tts.T(), block)
}

func (tts *TendermintTestSuite) TestQueryBlockResult() {
	result, err := tts.Tendermint().QueryBlockResult(1)
	require.NoError(tts.T(), err)
	require.Equal(tts.T(), int64(1), result.Height)
}

func (tts *TendermintTestSuite) TestQueryTx() {
	tx, err := tts.Tendermint().QueryTx("C270A774E11D466176F65AFB556387F340F485FDED411BFBC6EA219FEFDA4F56")
	require.NoError(tts.T(), err)
	fmt.Println(tx)

}

func (tts *TendermintTestSuite) TestQueryValidators() {
	result, err := tts.Tendermint().QueryValidators(1)
	require.NoError(tts.T(), err)
	require.Len(tts.T(), result.Validators, 1)
}

func (tts *TendermintTestSuite) TestQueryNodeInfo() {
	result, err := tts.Tendermint().QueryNodeInfo()
	require.NoError(tts.T(), err)
	fmt.Println(result)
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
