package service_test

import (
	"fmt"
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
	sts.NoError(err)
	sts.NotEmpty(result.Hash)

	defi, err := sts.Service().QueryDefinition(definition.ServiceName)
	sts.NoError(err)
	sts.Equal(definition.ServiceName, defi.Name)
	sts.Equal(definition.Description, defi.Description)
	sts.EqualValues(definition.Tags, defi.Tags)
	sts.Equal(definition.AuthorDescription, defi.AuthorDescription)
	sts.Equal(definition.Schemas, defi.Schemas)
	sts.Equal(sts.Account().Address, defi.Author)

	deposit, e := sdk.ParseDecCoins("20000iris")
	sts.NoError(e)
	binding := rpc.ServiceBindingRequest{
		ServiceName: definition.ServiceName,
		Deposit:     deposit,
		Pricing:     pricing,
	}
	result, err = sts.Service().BindService(binding, baseTx)
	sts.NoError(err)
	sts.NotEmpty(result.Hash)

	bindResp, err := sts.Service().QueryBinding(definition.ServiceName, sts.Account().Address)
	sts.NoError(err)
	sts.Equal(binding.ServiceName, bindResp.ServiceName)
	sts.Equal(sts.Account().Address, bindResp.Provider)
	sts.Equal(binding.Pricing, bindResp.Pricing)

	input := `{"pair":"iris-usdt"}`
	output := `{"last":"1:100"}`
	testResult := `{"code":200,"message":""}`

	var sub1 sdk.Subscription
	//sub1, err = sts.Service().RegisterSingleServiceListener(definition.ServiceName,
	//	func(reqID,input string) (string, string) {
	//		sts.Info().
	//          Str("reqID", reqID).
	//			Str("input", input).
	//			Str("output", output).
	//			Msg("provider received request")
	//		return output, testResult
	//	}, baseTx)

	router := rpc.ServiceRouter{
		definition.ServiceName: func(reqID, input string) (string, string) {
			sts.Info().
				Str("reqID", reqID).
				Str("input", input).
				Str("output", output).
				Msg("provider received request")
			return output, testResult
		},
	}
	sub1, err = sts.Service().RegisterServiceListener(router, baseTx)

	sts.NoError(err)

	serviceFeeCap, e := sdk.ParseDecCoins("1iris")
	sts.NoError(e)

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
	sub2, err = sts.Service().RegisterServiceResponseListener(requestContextID, func(reqCtxID string, responses string) {
		sts.Equal(reqCtxID, requestContextID)
		sts.Equal(output, responses)
		sts.Info().
			Str("requestContextID", requestContextID).
			Str("response", responses).
			Msg("consumer received response")
		sts.NoError(err)
		exit <- 1
	})

	sts.NoError(err)

	request, err := sts.Service().QueryRequestContext(requestContextID)
	sts.NoError(err)
	sts.Equal(request.ServiceName, invocation.ServiceName)
	sts.Equal(request.Input, invocation.Input)

	<-exit
	err = sts.Unsubscribe(sub1)
	err = sts.Unsubscribe(sub2)
	sts.NoError(err)
}

func (sts *ServiceTestSuite) TestQueryRequestContext() {
	reqCtxID := "E0F60DDA1140D90C1EEA4FD3D2C12570C042D346D80D288818A59103A8E3C6360000000000000000"
	_, err := sts.Service().QueryRequestContext(reqCtxID)
	sts.NoError(err)
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}
