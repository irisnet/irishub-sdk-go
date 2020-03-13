// Package oracle combines service to achieve decentralized injection from trusted Oracles such as Chainlink to IRISHub Oracle.
// Each data collection task is called a feed, and its underlying implementation depends on the service module.
//
// Feed life cycle It is basically the same as the Service RequestContext (paused, running), and is used to store off-chain data on the chain through Oracle nodes.
// In addition, you can only operate Feed through your Profiler account, and you cannot delete it.
//
// [More Details](https://www.irisnet.org/docs/features/oracle.html)
//
// As a quick start:
//
//
//	input := `{"pair":"iris-usdt"}`
//	feedName := generateFeedName(ots.serviceName)
//	serviceFeeCap, _ := sdk.ParseCoins("1000000000000000000iris-atto")
//
//	createFeedReq := rpc.FeedCreateRequest{
//		BaseTx:            ots.baseTx,
//		FeedName:          feedName,
//		LatestHistory:     5,
//		Description:       "fetch USDT-CNY ",
//		ServiceName:       ots.serviceName,
//		Providers:         []string{ots.Sender().String()},
//		Input:             input,
//		Timeout:           3,
//		ServiceFeeCap:     serviceFeeCap,
//		RepeatedFrequency: 5,
//		RepeatedTotal:     2,
//		AggregateFunc:     "avg",
//		ValueJsonPath:     "last",
//		ResponseThreshold: 1,
//	}
//	result, err := ots.Oracle().CreateFeed(createFeedReq)
//	require.NoError(ots.T(), err)
//	require.NotEmpty(ots.T(), result.Hash)
//
//	_, err = ots.Oracle().QueryFeed(feedName)
//	require.NoError(ots.T(), err)
//
//	result, err = ots.Oracle().StartFeed(feedName, ots.baseTx)
//	require.NoError(ots.T(), err)
//	require.NotEmpty(ots.T(), result.Hash)
//
//	for {
//		result, err := ots.Oracle().QueryFeedValue(feedName)
//		require.NoError(ots.T(), err)
//		if len(result) == int(createFeedReq.RepeatedTotal) {
//			goto stop
//		}
//	}
//
package oracle
