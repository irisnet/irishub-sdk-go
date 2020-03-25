package service_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"
	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	test.MockClient
	*log.Logger
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (sts *ServiceTestSuite) SetupTest() {
	sts.MockClient = test.NewMockClient()
	sts.Logger = log.NewLogger("info")
}

func (sts *ServiceTestSuite) TestService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":"1iris"}`

	baseTx := sdk.BaseTx{
		From:     sts.Account().Name,
		Gas:      20000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: sts.Account().Password,
	}

	definition := rpc.ServiceDefinitionRequest{
		ServiceName:       generateServiceName(),
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := sts.Service().DefineService(definition, baseTx)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), result.Hash)

	defi, err := sts.Service().QueryDefinition(definition.ServiceName)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), definition.ServiceName, defi.Name)
	require.Equal(sts.T(), definition.Description, defi.Description)
	require.EqualValues(sts.T(), definition.Tags, defi.Tags)
	require.Equal(sts.T(), definition.AuthorDescription, defi.AuthorDescription)
	require.Equal(sts.T(), definition.Schemas, defi.Schemas)
	require.Equal(sts.T(), sts.Account().Address, defi.Author)

	deposit, e := sdk.ParseDecCoins("20000iris")
	require.NoError(sts.T(), e)
	binding := rpc.ServiceBindingRequest{
		ServiceName: definition.ServiceName,
		Deposit:     deposit,
		Pricing:     pricing,
	}
	result, err = sts.Service().BindService(binding, baseTx)
	require.NoError(sts.T(), err)
	require.NotEmpty(sts.T(), result.Hash)

	bindResp, err := sts.Service().QueryBinding(definition.ServiceName, sts.Account().Address)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), binding.ServiceName, bindResp.ServiceName)
	require.Equal(sts.T(), sts.Account().Address, bindResp.Provider)
	require.Equal(sts.T(), binding.Pricing, bindResp.Pricing)

	input := `{"pair":"iris-usdt"}`
	output := `{"last":"1:100"}`
	testResult := `{"code":200,"message":""}`

	var sub1 sdk.Subscription
	router := rpc.ServiceRouter{
		definition.ServiceName: func(reqCtxID, reqID, input string) (string, string) {
			sts.Info().
				Str("reqCtxID", reqCtxID).
				Str("reqID", reqID).
				Str("input", input).
				Str("output", output).
				Msg("provider received request")
			_, err := sts.Service().QueryResponse(reqID)
			require.NoError(sts.T(), err)
			return output, testResult
		},
	}
	sub1, err = sts.Service().RegisterServiceRequestListener(router, baseTx)
	require.NoError(sts.T(), err)

	serviceFeeCap, e := sdk.ParseDecCoins("1iris")
	require.NoError(sts.T(), e)

	invocation := rpc.ServiceInvocationRequest{
		ServiceName:       definition.ServiceName,
		Providers:         []string{sts.Account().Address.String()},
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           1,
		SuperMode:         false,
		Repeated:          false,
		RepeatedFrequency: 3,
		RepeatedTotal:     1,
	}

	var requestContextID string
	var sub2 sdk.Subscription
	var exit = make(chan int, 0)

	requestContextID, err = sts.Service().InvokeService(invocation, baseTx)
	sub2, err = sts.Service().RegisterServiceResponseListener(requestContextID, func(reqCtxID, reqID, responses string) {
		sts.Info().
			Str("reqCtxID", reqCtxID).
			Str("reqID", reqID).
			Str("response", responses).
			Msg("consumer received response")

		require.Equal(sts.T(), reqCtxID, requestContextID)
		require.Equal(sts.T(), output, responses)
		request, err := sts.Service().QueryRequest(reqID)
		require.NoError(sts.T(), err)
		require.Equal(sts.T(), reqCtxID, request.RequestContextID)
		require.Equal(sts.T(), reqID, request.ID)
		require.Equal(sts.T(), input, request.Input)

		exit <- 1
	})

	require.NoError(sts.T(), err)

	request, err := sts.Service().QueryRequestContext(requestContextID)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), request.ServiceName, invocation.ServiceName)
	require.Equal(sts.T(), request.Input, invocation.Input)

	for {
		select {
		case <-exit:
			err = sts.Unsubscribe(sub1)
			err = sts.Unsubscribe(sub2)
			require.NoError(sts.T(), err)
		case <-time.After(1 * time.Minute):
			require.Panics(sts.T(), func() {}, "test service timeout")
		}
	}
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}
