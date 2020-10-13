package random_test

import (
	"github.com/irisnet/irishub-sdk-go/types/original"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/stretchr/testify/suite"
)

type RandomTestSuite struct {
	suite.Suite
	*test.MockClient
}

func TestSlashingTestSuite(t *testing.T) {
	suite.Run(t, new(RandomTestSuite))
}

func (rts *RandomTestSuite) SetupTest() {
	tc := test.GetMock()
	rts.MockClient = tc
}

func (rts *RandomTestSuite) TestGenerate() {
	baseTx := original.BaseTx{
		From:     rts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     original.Commit,
		Password: rts.Account().Password,
	}

	var memory = make(map[string]string, 1)
	var signal = make(chan int, 0)
	request := rpc.RandomRequest{
		BlockInterval: 2,
		Callback: func(reqID, randomNum string, err original.Error) {
			require.NoError(rts.T(), err)
			require.NoError(rts.T(), err)
			memory[reqID] = randomNum
			signal <- 1
		},
		Oracle: false,
	}

	reqID, err := rts.Random().Request(request, baseTx)
	require.NoError(rts.T(), err)
	memory[reqID] = ""
	<-signal
	require.NotEmpty(rts.T(), memory)
}
