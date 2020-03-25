package tendermint_test

import (
	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TendermintTestSuite struct {
	suite.Suite
	test.MockClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TendermintTestSuite))
}

func (tts *TendermintTestSuite) SetupTest() {
	tc := test.NewMockClient()
	tts.MockClient = tc
}

func (tts *TendermintTestSuite) TestQueryQueryBlock() {
	block, err := tts.Tendermint().QueryBlock(1)
	tts.NoError(err)
	tts.Equal(int64(1), block.Height)
}

func (tts *TendermintTestSuite) TestQueryBlockResult() {
	result, err := tts.Tendermint().QueryBlockResult(1)
	tts.NoError(err)
	tts.Equal(int64(1), result.Height)
}

func (tts *TendermintTestSuite) TestQueryTx() {
	coins, err := sdk.ParseDecCoins("0.1iris")
	tts.NoError(err)
	to := "faa1hp29kuh22vpjjlnctmyml5s75evsnsd8r4x0mm"
	baseTx := sdk.BaseTx{
		From:     tts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: tts.Account().Password,
	}

	result, err := tts.Bank().Send(to, coins, baseTx)
	tts.NoError(err)
	tts.NotEmpty(result.Hash)

	tx, err := tts.Tendermint().QueryTx(result.Hash)
	tts.NoError(err)
	tts.Equal(tx.Height, result.Height)
	tts.Equal(tx.Hash, result.Hash)

	builder := sdk.NewEventQueryBuilder().
		AddCondition(sdk.Cond(sdk.EventKey("tx.hash")).EQ(sdk.EventValue(result.Hash)))
	txs, err := tts.Tendermint().SearchTxs(builder, 1, 10)
	tts.NoError(err)
	tts.Equal(1, txs.Total)
	tts.Len(txs.Txs, 1)
	tts.Equal(result.Hash, txs.Txs[0].Hash)
}

func (tts *TendermintTestSuite) TestQueryValidators() {
	result, err := tts.Tendermint().QueryValidators(1)
	tts.NoError(err)
	tts.Len(result.Validators, 1)
}
