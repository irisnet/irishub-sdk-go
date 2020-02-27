package service_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/sim"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type ServiceTestSuite struct {
	suite.Suite
	sim.TestClient
	log.Logger
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (sts *ServiceTestSuite) SetupTest() {
	sts.TestClient = sim.NewClient()
	sts.Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

func (sts *ServiceTestSuite) TestService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":[{"denom":"iris-atto","amount":"1000000000000000000"}]}`

	baseTx := sdk.BaseTx{
		From: "test1",
		Gas:  20000,
		Fee:  "600000000000000000iris-atto",
		Memo: "test",
		Mode: sdk.Commit,
	}

	definition := sdk.ServiceDefinitionRequest{
		BaseTx:            baseTx,
		ServiceName:       generateServiceName(),
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := sts.DefineService(definition)
	require.NoError(sts.T(), err)
	require.True(sts.T(), result.IsSuccess())

	defi, err := sts.QueryDefinition(definition.ServiceName)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), definition.ServiceName, defi.Name)
	require.Equal(sts.T(), definition.Description, defi.Description)
	require.EqualValues(sts.T(), definition.Tags, defi.Tags)
	require.Equal(sts.T(), definition.AuthorDescription, defi.AuthorDescription)
	require.Equal(sts.T(), definition.Schemas, defi.Schemas)
	require.Equal(sts.T(), sts.Sender(), defi.Author)

	deposit, _ := sdk.ParseCoins("20000000000000000000000iris-atto")
	binding := sdk.ServiceBindingRequest{
		BaseTx:       baseTx,
		ServiceName:  definition.ServiceName,
		Deposit:      deposit,
		Pricing:      pricing,
		WithdrawAddr: "",
	}
	result, err = sts.BindService(binding)
	require.NoError(sts.T(), err)
	require.True(sts.T(), result.IsSuccess())

	bindResp, err := sts.QueryBinding(definition.ServiceName, sts.Sender())
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), binding.ServiceName, bindResp.ServiceName)
	require.Equal(sts.T(), sts.Sender(), bindResp.Provider)
	require.Equal(sts.T(), binding.Deposit.String(), bindResp.Deposit.String())
	require.Equal(sts.T(), binding.Pricing, bindResp.Pricing)

	input := `{"pair":"iris-usdt"}`
	output := `{"last":"1:100"}`

	err = sts.RegisterSingleInvocationListener(definition.ServiceName,
		func(input string) (string, string) {
			sts.Info("Provider received request", "input", input, "output", output)
			return output, ""
		}, baseTx)
	require.NoError(sts.T(), err)

	serviceFeeCap, _ := sdk.ParseCoins("1000000000000000000iris-atto")
	invocation := sdk.ServiceInvocationRequest{
		BaseTx:            baseTx,
		ServiceName:       definition.ServiceName,
		Providers:         []string{sts.Sender().String()},
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           3,
		SuperMode:         false,
		Repeated:          true,
		RepeatedFrequency: 5,
		RepeatedTotal:     -1,
	}
	var requestContextID string
	requestContextID, err = sts.InvokeService(invocation, func(reqCtxID string, response string) {
		require.Equal(sts.T(), reqCtxID, requestContextID)
		require.Equal(sts.T(), output, response)
		sts.Info("Consumer received response", "RequestContextID", requestContextID, "response", response)
	})

	sts.Info("Request success", "RequestContextID", requestContextID)
	require.NoError(sts.T(), err)

	request, err := sts.QueryRequestContext(requestContextID)
	require.NoError(sts.T(), err)
	require.Equal(sts.T(), request.ServiceName, invocation.ServiceName)
	require.Equal(sts.T(), request.Input, invocation.Input)

	time.Sleep(10 * time.Minute)
}

func (sts *ServiceTestSuite) TestQueryDefinition() {
	serviceName := "fbfbbc12-5872-11ea-a1dc-186590e06183"
	definition, err := sts.QueryDefinition(serviceName)
	require.NoError(sts.T(), err)
	fmt.Println(definition)
}

func (sts *ServiceTestSuite) TestQueryRequestContext() {
	reqCtxID := "0ef2acf4e1002f1e38c7f0df36be2ad11250aeb6343c6c7d9294f0424deecd96"
	definition, err := sts.QueryRequestContext(reqCtxID)
	require.NoError(sts.T(), err)
	fmt.Println(definition)
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}
