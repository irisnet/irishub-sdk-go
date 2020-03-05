package random_test

import (
	"testing"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RandomTestSuite struct {
	suite.Suite
	sim.TestClient
}

func TestSlashingTestSuite(t *testing.T) {
	suite.Run(t, new(RandomTestSuite))
}

func (rts *RandomTestSuite) SetupTest() {
	tc := sim.NewClient()
	rts.TestClient = tc
}

func (rts *RandomTestSuite) TestGenerate() {
	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	var memory = make(map[string]string, 1)
	var signal = make(chan int, 0)
	request := sdk.RandomRequest{
		BaseTx:        baseTx,
		BlockInterval: 2,
		Callback: func(reqID, randomNum string, err error) {
			require.NoError(rts.T(), err)
			memory[reqID] = randomNum
			signal <- 1
		},
	}
	reqID, err := rts.Random().Generate(request)
	require.NoError(rts.T(), err)
	memory[reqID] = ""
	<-signal
	require.NotEmpty(rts.T(), memory[reqID])
}
