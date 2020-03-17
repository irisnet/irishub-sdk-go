package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/irisnet/irishub-sdk-go/tools/log"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
	sdk "github.com/irisnet/irishub-sdk-go/types"
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
	pricing := `{"price":[{"denom":"iris-atto","amount":"1000000000000000000"}]}`

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

	err = sts.Service().RegisterSingleServiceListener(definition.ServiceName,
		func(input string) (string, string) {
			sts.Info().
				Str("input", input).
				Str("output", output).
				Msg("provider received request")
			return output, ""
		}, baseTx)
	require.NoError(sts.T(), err)

	serviceFeeCap, e := sdk.ParseDecCoins("1iris")
	require.NoError(sts.T(), e)

	invocation := rpc.ServiceInvocationRequest{
		ServiceName:       definition.ServiceName,
		Providers:         []string{sts.Account().Address.String()},
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           3,
		SuperMode:         false,
		Repeated:          false,
		RepeatedFrequency: 5,
		RepeatedTotal:     1,
	}
	var requestContextID string
	var exit = make(chan int, 0)
	requestContextID, err = sts.Service().InvokeService(invocation, func(reqCtxID string, response string) {
		require.Equal(sts.T(), reqCtxID, requestContextID)
		require.Equal(sts.T(), output, response)
		sts.Info().
			Str("requestContextID", requestContextID).
			Str("response", response).
			Msg("consumer received response")
		exit <- 1
	}, baseTx)

	sts.Info().
		Str("requestContextID", requestContextID).
		Msg("ServiceRequest service success")
	require.NoError(sts.T(), err)

	request, err := sts.Service().QueryRequestContext(requestContextID)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), request.ServiceName, invocation.ServiceName)
	require.Equal(sts.T(), request.Input, invocation.Input)

	<-exit
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}
