package random_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type RandomTestSuite struct {
	suite.Suite
	test.TestClient
}

func TestSlashingTestSuite(t *testing.T) {
	suite.Run(t, new(RandomTestSuite))
}

func (rts *RandomTestSuite) SetupTest() {
	tc := test.NewClient()
	rts.TestClient = tc
}

func (rts *RandomTestSuite) TestGenerate() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Memo: "test",
		Mode: sdk.Commit,
	}

	var memory = make(map[string]string, 1)
	var signal = make(chan int, 0)
	var callback = func(reqID, randomNum string, err sdk.Error) {
		require.NoError(rts.T(), err)
		memory[reqID] = randomNum
		signal <- 1
	}
	reqID, err := rts.Random().Request(2, callback, baseTx)
	require.NoError(rts.T(), err)
	memory[reqID] = ""
	<-signal
	require.NotEmpty(rts.T(), memory[reqID])
}
