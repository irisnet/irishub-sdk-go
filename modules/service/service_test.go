package service_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/utils/log"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	*test.MockClient
	*log.Logger
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (sts *ServiceTestSuite) SetupTest() {
	sts.MockClient = test.GetMock()
	sts.Logger = log.NewLogger("info")
}

//func (sts *ServiceTestSuite) TestService() {
//	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
//	pricing := `{"price":"1iris"}`
//
//	baseTx := sdk.BaseTx{
//		From:     sts.Account().Name,
//		Gas:      20000,
//		Memo:     "test",
//		Mode:     sdk.Commit,
//		Password: sts.Account().Password,
//	}
//
//	definition := rpc.ServiceDefinitionRequest{
//		ServiceName:       sts.RandStringOfLength(5),
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
//	require.Equal(sts.T(), sts.Account().Address, defi.Author)
//
//	deposit, e := sdk.ParseDecCoins("20000iris")
//	require.NoError(sts.T(), e)
//	binding := rpc.ServiceBindingRequest{
//		ServiceName: definition.ServiceName,
//		Deposit:     deposit,
//		Pricing:     pricing,
//		MinRespTime: 1,
//	}
//	result, err = sts.Service().BindService(binding, baseTx)
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), result.Hash)
//
//	bindResp, err := sts.Service().QueryBinding(definition.ServiceName, sts.Account().Address)
//	require.NoError(sts.T(), err)
//	require.Equal(sts.T(), binding.ServiceName, bindResp.ServiceName)
//	require.Equal(sts.T(), sts.Account().Address, bindResp.Provider)
//	require.Equal(sts.T(), binding.Pricing, bindResp.Pricing)
//
//	input := `{"pair":"iris-usdt"}`
//	output := `{"last":"1:100"}`
//	testResult := `{"code":200,"message":""}`
//
//	var sub1 sdk.Subscription
//	router := rpc.ServiceRegistry{
//		definition.ServiceName: func(reqCtxID, reqID, input string) (string, string) {
//			sts.Info().
//				Str("reqCtxID", reqCtxID).
//				Str("reqID", reqID).
//				Str("input", input).
//				Str("output", output).
//				Msg("provider received request")
//			_, err := sts.Service().QueryResponse(reqID)
//			require.NoError(sts.T(), err)
//			return output, testResult
//		},
//	}
//	sub1, err = sts.Service().SubscribeServiceRequest(router, baseTx)
//	require.NoError(sts.T(), err)
//
//	serviceFeeCap, e := sdk.ParseDecCoins("200iris")
//	require.NoError(sts.T(), e)
//
//	invocation := rpc.ServiceInvocationRequest{
//		ServiceName:       definition.ServiceName,
//		Providers:         []string{sts.Account().Address.String()},
//		Input:             input,
//		ServiceFeeCap:     serviceFeeCap,
//		Timeout:           1,
//		SuperMode:         false,
//		Repeated:          true,
//		RepeatedFrequency: 3,
//		RepeatedTotal:     -1,
//	}
//
//	var requestContextID string
//	var sub2 sdk.Subscription
//	var exit = make(chan int, 0)
//
//	requestContextID, err = sts.Service().InvokeService(invocation, baseTx)
//	require.NoError(sts.T(), err)
//
//	sub2, err = sts.Service().SubscribeServiceResponse(requestContextID, func(reqCtxID, reqID, responses string) {
//		sts.Info().
//			Str("reqCtxID", reqCtxID).
//			Str("reqID", reqID).
//			Str("response", responses).
//			Msg("consumer received response")
//
//		require.Equal(sts.T(), reqCtxID, requestContextID)
//		require.Equal(sts.T(), output, responses)
//		request, err := sts.Service().QueryRequest(reqID)
//		require.NoError(sts.T(), err)
//		require.Equal(sts.T(), reqCtxID, request.RequestContextID)
//		require.Equal(sts.T(), reqID, request.ID)
//		require.Equal(sts.T(), input, request.Input)
//
//		exit <- 1
//	})
//	require.NoError(sts.T(), err)
//
//	for {
//		select {
//		case <-exit:
//			err = sts.Unsubscribe(sub1)
//			err = sts.Unsubscribe(sub2)
//			require.NoError(sts.T(), err)
//			goto loop
//		case <-time.After(1 * time.Minute):
//			require.Panics(sts.T(), func() {}, "test service timeout")
//		}
//	}
//
//loop:
//	_, err = sts.Service().PauseRequestContext(requestContextID, baseTx)
//	require.NoError(sts.T(), err)
//
//	_, err = sts.Service().StartRequestContext(requestContextID, baseTx)
//	require.NoError(sts.T(), err)
//
//	request, err := sts.Service().QueryRequestContext(requestContextID)
//	require.NoError(sts.T(), err)
//	require.Equal(sts.T(), request.ServiceName, invocation.ServiceName)
//	require.Equal(sts.T(), request.Input, invocation.Input)
//
//	_, err = sts.Service().WithdrawEarnedFees(baseTx)
//	require.NoError(sts.T(), err)
//
//	addr, _, err := sts.Keys().Add(sts.RandStringOfLength(30), "1234567890")
//	require.NoError(sts.T(), err)
//	require.NotEmpty(sts.T(), addr)
//
//	_, err = sts.Service().SetWithdrawAddress(addr, baseTx)
//	require.NoError(sts.T(), err)
//}

func (sts *ServiceTestSuite) TestDefineService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`

	baseTx := sdk.BaseTx{
		From:     sts.Account().Name,
		Gas:      1000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: sts.Account().Password,
	}

	definition := rpc.ServiceDefinitionRequest{
		ServiceName:       "june30",
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}
	result, err := sts.Service().DefineService(definition, baseTx)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), result.Hash)
}

func (sts *ServiceTestSuite) TestQueryDefinition() {
	bindings, err := sts.Service().QueryDefinition("assettransfer")
	fmt.Println(bindings)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), bindings)
}

func (sts *ServiceTestSuite) TestQueryServiceBinding() {
	address, _ := sdk.AccAddressFromBech32("iaa13rtezlhpqms02syv27zc0lqc5nt3z4lcxzd9js")
	binding, err := sts.Service().QueryBinding("assettransfer", address)
	fmt.Println(binding)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), binding)
}

func (sts *ServiceTestSuite) TestQueryBindings() {
	bindings, err := sts.Service().QueryBindings("assettransfer")
	fmt.Println(bindings)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), bindings)
}

func (sts *ServiceTestSuite) TestQueryServiceRequest() {
	request, err := sts.Service().QueryRequest("2481D45863759DA4214546ABDE8B90C3FFCAE74DAE3D00A9C41785F4B04A17FA00000000000000000000000000000001000000000000BD0B0000")
	fmt.Println(request)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), request)
}

func (sts *ServiceTestSuite) TestQueryServiceRequests() {
	address, _ := sdk.AccAddressFromBech32("iaa13rtezlhpqms02syv27zc0lqc5nt3z4lcxzd9js")
	requests, err := sts.Service().QueryRequests("assettransfer", address)
	fmt.Println(requests)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), requests)
}

func (sts *ServiceTestSuite) TestQueryRequestsByReqCtx() {
	context, err := sts.Service().QueryRequestsByReqCtx("3A944731CCF66A297F52EC4C5B2E7975A1791703B6995952EE079F131A66762D0000000000000000", 1)
	fmt.Println(context)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), context)
}

func (sts *ServiceTestSuite) TestQueryResponse() {
	response, err := sts.Service().QueryResponse("2A639C624919577827812C3B83D5AE5C6E0941BD16AE5767D07A57A4434ABA2F00000000000000000000000000000001000000000000BF7E0000")
	fmt.Println(response)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), response)
}

func (sts *ServiceTestSuite) TestQueryResponses() {
	response, err := sts.Service().QueryResponses("1EE3DEC7A8F130C0B7EF7DD938BAA594BF9B5F161F2EB8A0097D83CFFFFD2D670000000000000000", 1)
	fmt.Println(response)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), response)
}

func (sts *ServiceTestSuite) TestQueryRequestContext() {
	context, err := sts.Service().QueryRequestContext("9BD523DAEA234A4D9779972630E0C48F7FC0B27D19639DEA37BE5B92C4C12E7C0000000000000000")
	fmt.Println(context)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), context)
}

func (sts *ServiceTestSuite) TestQueryFees() {
	fees, err := sts.Service().QueryFees("iaa13rtezlhpqms02syv27zc0lqc5nt3z4lcxzd9js")
	fmt.Println(fees)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), fees)
}
