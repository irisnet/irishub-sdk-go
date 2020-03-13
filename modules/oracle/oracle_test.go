package oracle_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/irisnet/irishub-sdk-go/rpc"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/irisnet/irishub-sdk-go/test"
	"github.com/irisnet/irishub-sdk-go/tools/log"
	sdk "github.com/irisnet/irishub-sdk-go/types"
)

type OracleTestSuite struct {
	suite.Suite
	test.TestClient
	*log.Logger
	serviceName string
	baseTx      sdk.BaseTx
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (ots *OracleTestSuite) SetupTest() {
	ots.TestClient = test.NewClient()
	ots.Logger = log.NewLogger("info")
}

func (ots *OracleTestSuite) SetupService() {
	schemas := `{"input":{"type":"object"},"output":{"type":"object"},"error":{"type":"object"}}`
	pricing := `{"price":[{"denom":"iris-atto","amount":"1000000000000000000"}]}`
	output := `{"last":"100"}`
	serviceName := generateServiceName()

	baseTx := sdk.BaseTx{
		From:     "test1",
		Gas:      20000,
		Memo:     "test",
		Mode:     sdk.Commit,
		Password: ots.Password(),
	}

	definition := rpc.ServiceDefinitionRequest{
		ServiceName:       serviceName,
		Description:       "this is a test service",
		Tags:              nil,
		AuthorDescription: "service provider",
		Schemas:           schemas,
	}

	result, err := ots.Service().DefineService(definition, baseTx)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	deposit, _ := sdk.ParseCoins("20000000000000000000000iris-atto")
	binding := rpc.ServiceBindingRequest{
		ServiceName: definition.ServiceName,
		Deposit:     deposit,
		Pricing:     pricing,
	}
	result, err = ots.Service().BindService(binding, baseTx)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	err = ots.Service().RegisterSingleServiceListener(serviceName,
		func(input string) (string, string) {
			ots.Info().Str("input", input).
				Str("output", output).
				Msg("Service received request")
			return output, ""
		}, baseTx)

	require.NoError(ots.T(), err)

	ots.serviceName = serviceName
	ots.baseTx = baseTx
}

func (ots *OracleTestSuite) TestFeed() {
	//before
	ots.SetupService()

	input := `{"pair":"iris-usdt"}`
	feedName := generateFeedName(ots.serviceName)
	serviceFeeCap, _ := sdk.ParseCoins("1000000000000000000iris-atto")

	createFeedReq := rpc.FeedCreateRequest{
		BaseTx:            ots.baseTx,
		FeedName:          feedName,
		LatestHistory:     5,
		Description:       "fetch USDT-CNY ",
		ServiceName:       ots.serviceName,
		Providers:         []string{ots.Sender().String()},
		Input:             input,
		Timeout:           3,
		ServiceFeeCap:     serviceFeeCap,
		RepeatedFrequency: 5,
		RepeatedTotal:     2,
		AggregateFunc:     "avg",
		ValueJsonPath:     "last",
		ResponseThreshold: 1,
	}
	result, err := ots.Oracle().CreateFeed(createFeedReq)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	_, err = ots.Oracle().QueryFeed(feedName)
	require.NoError(ots.T(), err)

	result, err = ots.Oracle().StartFeed(feedName, ots.baseTx)
	require.NoError(ots.T(), err)
	require.NotEmpty(ots.T(), result.Hash)

	for {
		result, err := ots.Oracle().QueryFeedValue(feedName)
		require.NoError(ots.T(), err)
		if len(result) == int(createFeedReq.RepeatedTotal) {
			goto stop
		}
	}
stop:
	ots.Info().Msg("test feed success")
}

func generateServiceName() string {
	return fmt.Sprintf("service-%d", time.Now().Nanosecond())
}

func generateFeedName(serviceName string) string {
	return fmt.Sprintf("feed-%s", serviceName)
}
