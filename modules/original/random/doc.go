// Package random describes the usage and scope of random numbers on IRISHub.
// This feature is currently in beta and please assess the risk yourself before using.
//
// [More Details](https://www.irisnet.org/docs/features/random.html)
//
// As a quick start:
//
//	baseTx := sdk.BaseTx{
//		From: "test1",
//		Gas:  20000,
//		Memo: "test",
//		Mode: sdk.Commit,
//	}
//
//	var memory = make(map[string]string, 1)
//	var signal = make(chan int, 0)
//	request := rpc.RandomRequest{
//		BlockInterval: 2,
//		Callback: func(reqID, randomNum string, err sdk.Error) {
//			require.NoError(rts.T(), err)
//			require.NoError(rts.T(), err)
//			memory[reqID] = randomNum
//			signal <- 1
//		},
//		Oracle: false,
//	}
//
//	reqID, err := rts.Random().Request(request, baseTx)
//	require.NoError(rts.T(), err)
//	memory[reqID] = ""
//	<-signal
//	require.NotEmpty(rts.T(), memory[reqID])
//
package random
