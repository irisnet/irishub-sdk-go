// Package service bridge the gap between the blockchain world and the conventional business application world,
// by mediating a complete lifecycle of off-chain services -- from their definition, binding (provider registration), invocation, to their governance (profiling and dispute resolution).
//
// By enhancing the IBC processing logic to support service semantics, the IRIS SDK is intended to allow distributed business services to be available across the internet of blockchains.
// The Interface description language (IDL) we introduced is to work with the service standardized definitions to satisfy service invocations across different programming languages. The currently supported IDL language is protobuf
//
// [More Details](https://www.irisnet.org/docs/features/service.html)
//
// As a quick start:
//
//	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
//	pricing := `{"price":[{"denom":"iris-atto","amount":"1000000000000000000"}]}`
//
//	baseTx := sdk.BaseTx{
//		From: "test1",
//		Gas:  20000,
//		Memo: "test",
//		Mode: sdk.Commit,
//	}
//
//	definition := rpc.ServiceDefinitionRequest{
//		ServiceName:       generateServiceName(),
//		Description:       "this is a test service",
//		Tags:              nil,
//		AuthorDescription: "service provider",
//		Schemas:           schemas,
//	}
//
//	result, err := sts.Service().DefineService(definition, baseTx)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), result.Hash)
//
//	defi, err := sts.Service().QueryDefinition(definition.ServiceName)
//	require.NoError(sts.T(), err)
//	require.Equal(sts.T(), definition.ServiceName, defi.Name)
//	require.Equal(sts.T(), definition.Description, defi.Description)
//	require.EqualValues(sts.T(), definition.Tags, defi.Tags)
//	require.Equal(sts.T(), definition.AuthorDescription, defi.AuthorDescription)
//	require.Equal(sts.T(), definition.Schemas, defi.Schemas)
//	require.Equal(sts.T(), sts.Sender(), defi.Author)
//
//	deposit, _ := sdk.ParseCoins("20000000000000000000000iris-atto")
//	binding := rpc.ServiceBindingRequest{
//		ServiceName: definition.ServiceName,
//		Deposit:     deposit,
//		Pricing:     pricing,
//	}
//	result, err = sts.Service().BindService(binding, baseTx)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), result.Hash)
//
//	bindResp, err := sts.Service().QueryBinding(definition.ServiceName, sts.Sender())
//	require.NoError(sts.T(), err)
//	require.Equal(sts.T(), binding.ServiceName, bindResp.ServiceName)
//	require.Equal(sts.T(), sts.Sender(), bindResp.Provider)
//	require.Equal(sts.T(), binding.Deposit.String(), bindResp.Deposit.String())
//	require.Equal(sts.T(), binding.Pricing, bindResp.Pricing)
//
//	input := `{"pair":"iris-usdt"}`
//	output := `{"last":"1:100"}`
//
//	err = sts.Service().RegisterSingleServiceListener(definition.ServiceName,
//		func(input string) (string, string) {
//			sts.Info().
//				Str("input", input).
//				Str("output", output).
//				Msg("provider received request")
//			return output, ""
//		}, baseTx)
//	require.NoError(sts.T(), err)
//
//	serviceFeeCap, _ := sdk.ParseCoins("1000000000000000000iris-atto")
//	invocation := rpc.ServiceInvocationRequest{
//		ServiceName:       definition.ServiceName,
//		Providers:         []string{sts.Sender().String()},
//		Input:             input,
//		ServiceFeeCap:     serviceFeeCap,
//		Timeout:           3,
//		SuperMode:         false,
//		Repeated:          true,
//		RepeatedFrequency: 5,
//		RepeatedTotal:     -1,
//	}
//	var requestContextID string
//	var exit = make(chan int, 0)
//	requestContextID, err = sts.Service().InvokeService(invocation, func(reqCtxID string, response string) {
//		require.Equal(sts.T(), reqCtxID, requestContextID)
//		require.Equal(sts.T(), output, response)
//		sts.Info().
//			Str("requestContextID", requestContextID).
//			Str("response", response).
//			Msg("consumer received response")
//		exit <- 1
//	}, baseTx)
//
//	sts.Info().
//		Str("requestContextID", requestContextID).
//		Msg("ServiceRequest service success")
//	require.NoError(sts.T(), err)
//
//	request, err := sts.Service().QueryRequestContext(requestContextID)
//	require.NoError(sts.T(), err)
//	require.Equal(sts.T(), request.ServiceName, invocation.ServiceName)
//	require.Equal(sts.T(), request.Input, invocation.Input)
//
//	<-exit
//
package service
